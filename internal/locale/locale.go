package locale

import (
	waqi "air-quality-bot/pkg/waqi"
)

const (
	Greeting = iota
	AirQualityBtn
	AirQualityInfoMsg
	AirQualityLocationMsg
	AirQualityPollutionLvlBtn
	AirQualityPollutionLvlMsg
	LocationBtn
	NotificationsBtn
	NotificationsMsq
	NotificationsCreateBtn
	NotificationViewBtn
	NotificationLocationMsg
	NotificationCreatedMsg
	NotificationsUserListNotFoundMsg
	NotificationsUserListMsg
	NotificationDetailsBtn
	NotificationDetailsMsg
	NotificationDetailsPauseBtn
	NotificationDetailsUnpauseBtn
	NotificationDetailsEditTimeBtn
	NotificationDetailsEditTimeZoneBtn
	NotificationDetailsEditLocationBtn
	NotificationDetailsDeleteBtn
	NotificationStatusActiveMsg
	NotificationStatusPausedMsg
	NotificationDetailsManageMsg
	NotificationUnpausedMsg
	NotificationPausedMsg
	NotificationDeletedMsg
	NotificationEditLocationMsg
	NotificationLocationUpdatedMsg
	NotificationEditTimeMsg
	NotificationTimeUpdatedMsg
	NotificationEditTimeZoneMsg
	NotificationTimeZoneUpdatedMsg
	NotificationEditBackToDetailsBtn
	NotificationBackToListBtn
	CancelBtn
	ContinueBtn
	MenuMsg
	LocationProcessingMsg
	TimeZoneClarificationMsg
	TimeZoneRightBtn
	TimeZoneWrongBtn
	TimeZoneManualMsg
	TimeZoneUpdatedMsg
	TimePickerMsg
	PollutionLvlGood
	PollutionLvlModerate
	PollutionLvlSensUnhealthy
	PollutionLvlUnhealthy
	PollutionLvlVeryUnhealthy
	PollutionLvlHazardous
	ErrorMsg
	ErrorLoadLocationMsg
	AccessDenied
	MessageNotRecognisedMsg
)

func Get(msgId int, langCode string) string {
	switch langCode {
	default:
		return engMessages[msgId]
	}
}

const (
	AirQualityIndexScaleImg = iota
)

func GetImage(imgId int, langCode string) string {
	switch langCode {
	default:
		return engImages[imgId]
	}
}

func GetPollutionLvl(lvl waqi.PollutionLvl, langCode string) string {
	if lvl == waqi.PollutionLvlGood {
		return Get(PollutionLvlGood, langCode)
	}
	if lvl == waqi.PollutionLvlModerate {
		return Get(PollutionLvlModerate, langCode)
	}
	if lvl == waqi.PollutionLvlSensUnhealthy {
		return Get(PollutionLvlSensUnhealthy, langCode)
	}
	if lvl == waqi.PollutionLvlUnhealthy {
		return Get(PollutionLvlUnhealthy, langCode)
	}
	if lvl == waqi.PollutionLvlVeryUnhealthy {
		return Get(PollutionLvlVeryUnhealthy, langCode)
	}
	return Get(PollutionLvlHazardous, langCode)
}
