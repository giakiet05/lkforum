package model

type VNProvince string

const (
	ProvinceTuyenQuang VNProvince = "Tuyên Quang"
	ProvinceLaoCai     VNProvince = "Lào Cai"
	ProvinceThaiNguyen VNProvince = "Thái Nguyên"
	ProvincePhuTho     VNProvince = "Phú Thọ"
	ProvinceBacNinh    VNProvince = "Bắc Ninh"
	ProvinceHungYen    VNProvince = "Hưng Yên"
	ProvinceHaiPhong   VNProvince = "Hải Phòng"
	ProvinceNinhBinh   VNProvince = "Ninh Bình"
	ProvinceQuangTri   VNProvince = "Quảng Trị"
	ProvinceDaNang     VNProvince = "Đà Nẵng"
	ProvinceQuangNgai  VNProvince = "Quảng Ngãi"
	ProvinceGiaLai     VNProvince = "Gia Lai"
	ProvinceKhanhHoa   VNProvince = "Khánh Hòa"
	ProvinceDienBien   VNProvince = "Điện Biên"
	ProvinceHaNoi      VNProvince = "Hà Nội"
	ProvinceHaTinh     VNProvince = "Hà Tĩnh"
	ProvinceLangSon    VNProvince = "Lạng Sơn"
	ProvinceLaiChau    VNProvince = "Lai Châu"
	ProvinceNgheAn     VNProvince = "Nghệ An"
	ProvinceQuangNinh  VNProvince = "Quảng Ninh"
	ProvinceSonLa      VNProvince = "Sơn La"
	ProvinceThanhHoa   VNProvince = "Thanh Hóa"
	ProvinceCaoBang    VNProvince = "Cao Bằng"
	ProvinceHue        VNProvince = "TP. Huế"
	ProvinceLamDong    VNProvince = "Lâm Đồng"
	ProvinceDakLak     VNProvince = "Đắk Lắk"
	ProvinceHCM        VNProvince = "TPHCM"
	ProvinceDongNai    VNProvince = "Đồng Nai"
	ProvinceTayNinh    VNProvince = "Tây Ninh"
	ProvinceCanTho     VNProvince = "Cần Thơ"
	ProvinceVinhLong   VNProvince = "Vĩnh Long"
	ProvinceDongThap   VNProvince = "Đồng Tháp"
	ProvinceCaMau      VNProvince = "Cà Mau"
	ProvinceAnGiang    VNProvince = "An Giang"
)

var AllProvinces = []VNProvince{
	ProvinceTuyenQuang,
	ProvinceLaoCai,
	ProvinceThaiNguyen,
	ProvincePhuTho,
	ProvinceBacNinh,
	ProvinceHungYen,
	ProvinceHaiPhong,
	ProvinceNinhBinh,
	ProvinceQuangTri,
	ProvinceDaNang,
	ProvinceQuangNgai,
	ProvinceGiaLai,
	ProvinceKhanhHoa,
	ProvinceDienBien,
	ProvinceHaNoi,
	ProvinceHaTinh,
	ProvinceLangSon,
	ProvinceLaiChau,
	ProvinceNgheAn,
	ProvinceQuangNinh,
	ProvinceSonLa,
	ProvinceThanhHoa,
	ProvinceCaoBang,
	ProvinceHue,
	ProvinceLamDong,
	ProvinceDakLak,
	ProvinceHCM,
	ProvinceDongNai,
	ProvinceTayNinh,
	ProvinceCanTho,
	ProvinceVinhLong,
	ProvinceDongThap,
	ProvinceCaMau,
	ProvinceAnGiang,
}

func GetAllProvinces() []string {
	provinces := make([]string, len(AllProvinces))
	for i, p := range AllProvinces {
		provinces[i] = string(p)
	}
	return provinces
}

func IsValidProvince(province VNProvince) bool {
	for _, p := range AllProvinces {
		if p == province {
			return true
		}
	}
	return false
}

func (p VNProvince) String() string {
	return string(p)
}
