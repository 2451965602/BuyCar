package AIAgent

import (
	"buycar/config"
	"buycar/pkg/errno"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func CallOpenAICompat(ctx context.Context, model string, messages []Message) (string, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	reqBody := struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}{
		Model:    model,
		Messages: messages,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", config.AiEndpoint.Url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+config.AiEndpoint.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, "API 请求失败: "+err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, " API 请求失败，状态码: "+resp.Status+"，响应内容: "+string(b))
	}

	var result struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, "解析响应失败: "+err.Error())
	}
	if len(result.Choices) == 0 {
		return "", errno.NewErrNo(errno.InternalServiceErrorCode, "API 响应中没有可用的回答")
	}
	return result.Choices[0].Message.Content, nil
}
