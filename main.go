package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// StreamData sends events to the client
func StreamData(ctx *gin.Context) {
	// Get the ID from the URL
	tickerQuery := ctx.Query("ticker")

	log.Println("tickerQuery:", tickerQuery)

	// Set headers for SSE
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	// Flush the headers
	ctx.Writer.Flush()

	// CloseNotify returns a channel that gets closed when the client disconnects
	clientGone := ctx.Writer.CloseNotify()

	// Adjusted ticker to 1 second (for practical use cases)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			// Send the event in SSE format
			ctx.Writer.Write([]byte(fmt.Sprintf("id: %d\n", t.UnixMilli())))
			ctx.Writer.Write([]byte(fmt.Sprintf("event: %s\n", "message")))
			ctx.Writer.Write([]byte(fmt.Sprintf("data: %s\n", "Công ty Cổ phần Tập đoàn Hoà Phát (HPG) là một trong những Tập đoàn sản xuất công nghiệp đa ngành tại Việt Nam. Khởi đầu từ một Công ty chuyên buôn bán các loại máy xây dựng từ tháng 8/1992, hiện tại Tập đoàn Hòa Phát hoạt động chủ yếu trong các lĩnh vực gang thép, sản phẩm thép, điện máy gia dụng, nông nghiệp và bất động sản. Trong đó, lĩnh vực Thép  đóng vai trò chủ đạo và là mảng kinh doanh cốt lõi của tập đoàn với việc đóng góp hơn 90% doanh thu và lợi nhuận. HPG hiện là doanh nghiệp sản xuất thép xây dựng và ống thép lớn nhất Việt Nam với thị phần lần lượt là 36.4% và 29.07% (tháng 7/2022), và là doanh nghiệp Việt Nam duy nhất sản xuất được Thép cuộn cán nóng HRC. HPG được niêm yết và giao dịch Sở Giao dịch Chứng khoán Thành phố Hồ Chí Minh (HOSE) từ năm 2007.")))
			ctx.Writer.Write([]byte("\n"))
			ctx.Writer.Flush()

		case <-clientGone:
			// When the client closes the connection
			log.Println("Connection closed by client")
			return
		}
	}
}

func main() {
	r := gin.Default()

	// Access-Control-Allow-Origin
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// Define the SSE route
	r.GET("/events", StreamData)

	// Start the server
	log.Println("Starting server on :8080")
	r.Run(":8080")
}
