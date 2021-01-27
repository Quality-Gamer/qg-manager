package conf

import "strconv"

const Game = "gm"
const Manager = "mg"
const Model = "md"
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
const ManagerOccurrence = "mo"
const NumberOccurrences = "no"
const CurrentWeek = "cw"
const CurrentMoney = "cm"
const CurrentTime = "ct"
const Level = "lv"
const Process = "pc"
const Score = "sc"
const Resource = "rc"
const Activity = "at"
const None = ""
const Solve = "sv"
const Identifier = "id"
const InitialTime = "it"
const InitialMoney = "im"
const Type = "tp"
const BelongsTo = "bt"
const Time = "tm"
const Price = "pr"
const Quantity = "qt"
const Count = "ct"
const ResourcesIds = "rr"
const ActivityIds = "aa"
const Member = "mb"
const Product = "pd"
const Open = "O"
const Close = "C"

func GetKeyManager(userId,week int, managerId string) string {
	return Game + ":" + Manager + ":" + strconv.Itoa(userId) + ":" + managerId + ":" + strconv.Itoa(week)
}

func GetKeyOccurrence(userId int,managerId string) string {
	return Game + ":" + Manager + ":" + strconv.Itoa(userId) + ":" + managerId
}

//private
func getModelKey() string {
	return Game + ":" + Manager + ":" + Model
}

func GetGameModelKey(id,field string) string {
	return getModelKey() + ":" + id + ":" + field
}

func GetLevelKey(id,field string, level int) string {
	return GetGameModelKey(id,Level) + ":" + strconv.Itoa(level) + ":" + field
}

func GetProcessKey(id,field string,level,process int) string {
	return GetGameModelKey(id,Level) + ":" + strconv.Itoa(level) + ":" + Process  + ":" + strconv.Itoa(process) + ":" + field
}

func GetLevelProcessKey(id string,level,process int) string {
	return GetGameModelKey(id,Level) + ":" + strconv.Itoa(level) + ":" + Process  + ":" + strconv.Itoa(process)
}

func GetActivityKey(id,field string,level,process int,unit string) string {
	return GetLevelProcessKey(id,level,process) + ":" + Activity + ":" + unit + ":" + field
}

func GetResourceKey(id,field string,level,process int,unit string) string {
	return GetLevelProcessKey(id,level,process) + ":" + Resource + ":" + unit + ":" + field
}

func GetGameOccurrenceKey(id,field string) string {
	return GetGameModelKey(id,Occurrence) + ":" + field
}

func GetGameOccurrenceKeyType(id,_type,field string) string {
	return GetGameModelKey(id,Occurrence) + ":" + _type + ":" + field
}

func GetSolveKey(id,field string) string {
	return GetGameModelKey(id,Solve) + ":" + field
}

func GetScoreKey(id,field string) string {
	return GetGameModelKey(id,Score) + ":" + field
}
