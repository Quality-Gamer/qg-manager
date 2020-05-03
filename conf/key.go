package conf

import "strconv"

const Game = "gm"
const Manager = "mg"
const ChallengeId = "ci"
const UserId = "ui"
const ProjectId = "pi"
const Week = "wk"
const Progress = "pg"
const Team = "tm"
const RequirementAnalyst = "ra"
const ProductOwner = "po"
const Designer = "dg"
const Frontend = "ft"
const Tester = "tt"
const Backend = "bk"
const CustomerContact = "cc"
const Delivery = "dv"
const RiskAnalysis = "rk"
const Scrum = "sc"
const License = "lc"
const Ide = "li"
const DesignSoftware = "ld"
const Status = "st"
const Event = "ev"
const Name = "nm"
const Action = "ac"
const Occurrence = "oc"
const OccurrenceId = "oi"
const Description = "dc"
const UserOccurrence = "uo"
const NumberOccurrences = "no"
const CurrentWeek = "cw"
const CurrentMoney = "cm"
const CurrentTime = "ct"

func GetKeyManager(userId,week int, managerId string) string {
	return Game + ":" + Manager + ":" + strconv.Itoa(userId) + ":" + managerId + ":" + strconv.Itoa(week)
}

func GetKeyOccurrence(userId int,managerId string) string {
	return Game + ":" + Manager + ":" + strconv.Itoa(userId) + ":" + managerId
}