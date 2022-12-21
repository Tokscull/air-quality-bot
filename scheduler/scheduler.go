package scheduler

import (
	"air-quality-bot/internal/locale"
	"air-quality-bot/internal/notification"
	"air-quality-bot/pkg/logger"
	"air-quality-bot/pkg/utils"
	"air-quality-bot/pkg/waqi"
	tgBot "air-quality-bot/telegram/bot"
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"time"
)

const (
	NotificationsPageSize = 25
)

type Scheduler struct {
	log              *logger.Logger
	bot              *tgBot.Bot
	aqiService       *waqi.Service
	notificationRepo notification.Repository
}

func NewScheduler(log *logger.Logger, bot *tgBot.Bot, aqiService *waqi.Service, ntfRepo notification.Repository) *Scheduler {
	s := &Scheduler{
		log:              log,
		bot:              bot,
		aqiService:       aqiService,
		notificationRepo: ntfRepo,
	}
	return s
}

func (s Scheduler) Start() *gocron.Scheduler {
	sch := gocron.NewScheduler(time.UTC)

	//every hour
	_, _ = sch.Cron("0 */1 * * *").Do(s.notificationJob)

	//every minute - test
	//_, _ = sch.Cron("*/1 * * * *").Do(s.notificationJob)

	sch.StartAsync()
	s.log.Info("Scheduler started")
	return sch
}

func (s Scheduler) notificationJob() {
	now := time.Now().UTC()
	s.log.Info("Notification job started", zap.Time("now", now))

	totalPass, totalFail := s.sendNotifications(now, 1, NotificationsPageSize, 0, 0)

	s.log.Info("Notification job end",
		zap.Time("now", now),
		zap.Int("totalPass", totalPass),
		zap.Int("totalFail", totalFail),
	)
}

func (s Scheduler) sendNotifications(ntfTime time.Time, page int, pageSize int, totalPass int, totalFail int) (int, int) {
	offset := (page - 1) * pageSize
	rows, dbErr := s.notificationRepo.FindAllByNotifyAtPerPage(context.Background(), utils.TruncateToHour(ntfTime), pageSize, offset)
	if dbErr != nil || len(rows) == 0 {
		s.log.Warn("rows for notification job not found", zap.Error(dbErr))
		return totalPass, totalFail
	}

	for _, ntf := range rows {
		err := s.sendNotificationMessage(&ntf)
		if err != nil {
			totalFail += 1
		} else {
			totalPass += 1
			dbErr = s.notificationRepo.UpdateLastTimeProcessedAt(context.Background(), time.Now().UTC(), ntf.ID)
			if dbErr != nil {
				s.log.Error("Error updating notification last_time_processed_at", zap.Error(dbErr))
			}
		}
	}

	time.Sleep(3 * time.Second)
	return s.sendNotifications(ntfTime, page+1, pageSize, totalPass, totalFail)
}

func (s Scheduler) sendNotificationMessage(ntf *notification.Notification) error {
	resp, err := s.aqiService.GetByGeo(ntf.Location.Latitude, ntf.Location.Longitude)
	if err != nil {
		return err
	}

	textForm := locale.Get(locale.AirQualityInfoMsg, ntf.User.LangCode)
	pollutionLvl := locale.GetPollutionLvl(resp.Level, ntf.User.LangCode)

	text := fmt.Sprintf(textForm, resp.Station.Name, resp.AQI, pollutionLvl, resp.Station.URL)
	return s.bot.SendMessage(ntf.User.ChatID, text)
}
