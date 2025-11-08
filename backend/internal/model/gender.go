package model

type Gender string

const (
	GenderMale           Gender = "male"
	GenderFemale         Gender = "female"
	GenderPreferNotToSay Gender = "prefer_not_to_say"
)

var AllGenders = []Gender{
	GenderMale,
	GenderFemale,
	GenderPreferNotToSay,
}

func GetAllGenders() []string {
	genders := make([]string, len(AllGenders))
	for i, g := range AllGenders {
		genders[i] = string(g)
	}
	return genders
}

func IsValidGender(gender Gender) bool {
	for _, g := range AllGenders {
		if g == gender {
			return true
		}
	}
	return false
}

func (g Gender) String() string {
	return string(g)
}
