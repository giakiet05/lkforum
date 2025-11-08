package model

type Interest string

const (
	// Technology
	InterestProgramming Interest = "Lập trình"

	// Creative
	InterestDesign      Interest = "Thiết kế"
	InterestPhotography Interest = "Nhiếp ảnh"
	InterestVideo       Interest = "Dựng video"
	InterestMusic       Interest = "Âm nhạc"
	InterestWriting     Interest = "Viết lách"

	// Lifestyle
	InterestGaming  Interest = "Gaming"
	InterestSports  Interest = "Thể thao"
	InterestTravel  Interest = "Du lịch"
	InterestCooking Interest = "Nấu ăn"
	InterestFitness Interest = "Gym & Fitness"
	InterestBooks   Interest = "Đọc sách"

	// Business & Career
	InterestBusiness  Interest = "Kinh doanh"
	InterestMarketing Interest = "Marketing"
	InterestFinance   Interest = "Tài chính"
	InterestStartup   Interest = "Startup"

	// Education
	InterestStudy     Interest = "Học tập"
	InterestScience   Interest = "Khoa học"
	InterestLanguages Interest = "Ngoại ngữ"

	// Entertainment
	InterestMovies Interest = "Phim ảnh"
	InterestAnime  Interest = "Anime"
	InterestKpop   Interest = "K-pop"
	InterestPets   Interest = "Thú cưng"
)

var AllInterests = []Interest{
	InterestProgramming, InterestDesign, InterestPhotography,
	InterestVideo, InterestMusic, InterestWriting, InterestGaming,
	InterestSports, InterestTravel, InterestCooking, InterestFitness,
	InterestBooks, InterestBusiness, InterestMarketing, InterestFinance,
	InterestStartup, InterestStudy, InterestScience, InterestLanguages,
	InterestMovies, InterestAnime, InterestKpop, InterestPets,
}

func GetAllInterests() []string {
	interests := make([]string, len(AllInterests))
	for i, interest := range AllInterests {
		interests[i] = string(interest)
	}
	return interests
}

func IsValidInterest(interest Interest) bool {
	for _, i := range AllInterests {
		if i == interest {
			return true
		}
	}
	return false
}

func (i Interest) String() string {
	return string(i)
}
