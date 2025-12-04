package gemini

import (
	"bytes"
	"fmt"
	"text/template"
)

const moderationPromptTemplate = `Bạn là AI moderator cho diễn đàn cộng đồng.

**NHIỆM VỤ:** Phân tích nội dung và xác định có vi phạm tiêu chuẩn cộng đồng không.

**TIÊU CHUẨN CỘNG ĐỒNG:**
1. ❌ Hate Speech: Phân biệt chủng tộc, tôn giáo, giới tính, LGBTQ+, kỳ thị
2. ❌ Bạo lực: Đe dọa, kích động bạo lực, hình ảnh máu me, nội dung gây sốc
3. ❌ NSFW: Nội dung khiêu dâm, khỏa thân, tình dục
4. ❌ Spam: Quảng cáo thương mại, lừa đảo, scam, clickbait
5. ❌ Quấy rối: Tấn công cá nhân, doxxing, bullying, xúc phạm
6. ❌ Thông tin sai lệch: Tin giả nguy hiểm về y tế, chính trị

**NỘI DUNG KIỂM TRA:**
{{if .Title}}Tiêu đề: {{.Title}}
{{end}}{{if .Text}}Nội dung: {{.Text}}
{{end}}{{if .HasImages}}[Kèm {{.ImageCount}} ảnh - đang phân tích]
{{end}}{{if .HasVideos}}[Kèm {{.VideoCount}} video - đang phân tích thumbnail]
{{end}}

**YÊU CẦU TRẢ VỀ JSON:**
{
  "is_violation": boolean,
  "confidence": 0.0-1.0,
  "categories": ["hate_speech", "violence", "nsfw", "spam", "harassment", "misinformation"],
  "reason": "Giải thích NGẮN GỌN bằng tiếng Việt (tối đa 1-2 câu) tại sao vi phạm hoặc an toàn"
}

**LƯU Ý:**
- Nếu không chắc chắn (confidence < 0.7) → is_violation = false
- Chỉ reject nếu vi phạm RÕ RÀNG
- Bỏ qua lỗi chính tả, ngữ pháp
- Cho phép tranh luận, phê bình lịch sự
- Nếu nội dung trống → is_violation = false

Chỉ trả về JSON, không giải thích thêm.`

func buildModerationPrompt(req *ContentCheckRequest) string {
	tmpl := template.Must(template.New("prompt").Parse(moderationPromptTemplate))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, map[string]interface{}{
		"Title":      req.Title,
		"Text":       req.Text,
		"HasImages":  len(req.ImageURLs) > 0,
		"ImageCount": len(req.ImageURLs),
		"HasVideos":  len(req.VideoURLs) > 0,
		"VideoCount": len(req.VideoURLs),
	})

	if err != nil {
		return fmt.Sprintf("Phân tích nội dung này:\nTiêu đề: %s\nNội dung: %s", req.Title, req.Text)
	}

	return buf.String()
}
