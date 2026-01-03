package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/dylanewe/moni/internal/store"
	"github.com/openai/openai-go/v3"
)

type LLMParserService struct {
	client *openai.Client
}

func (s *LLMParserService) ParseStatement(ctx context.Context, categories []string, filepath string) ([]store.Transaction, error) {
	parsedData, err := extractPDFToText(ctx, filepath)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`You are a financial data extraction specialist. Analyze the provided financial document and extract transaction data following these rules:

DOCUMENT TYPES:
- Payslip: Extract ONLY the net salary (final take-home pay)
- Bank/Credit statements: Extract all transactions

EXTRACTION RULES:
1. Amount signs:
   - Negative (-) for: expenses, purchases, withdrawals, fees, charges
   - Positive (+) for: income, salary, deposits, refunds, interest earned, cashback
2. Skip these transactions entirely:
   - Credit card payments from bank accounts
   - Credit card settlement transactions
   - Balance transfers between own accounts
   - Duplicate entries
3. Category: Match to the provided list. If uncertain or no clear match, use empty string ""
4. Date: Convert all dates to YYYY-MM-DD format. If year is ambiguous, infer from context
5. Description: Use the merchant/payee name or transaction description as-is

AVAILABLE CATEGORIES:
%s

OUTPUT FORMAT:
Return a valid JSON array of transactions. Each transaction must follow this exact structure:
{
  "description": "string - merchant or transaction description",
  "category": "string - from available categories or empty",
  "amount": float - negative for expenses, positive for income,
  "date": "string - YYYY-MM-DD format"
}

Return ONLY the JSON array with no additional text or markdown formatting.

DOCUMENT TEXT:
`, strings.Join(categories, ", "))

	chatCompletion, err := s.client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
			openai.UserMessage(parsedData),
		},
		Model:       openai.ChatModelGPT4oMini,
		Temperature: openai.Float(0),
	})
	if err != nil {
		return nil, err
	}

	content := chatCompletion.Choices[0].Message.Content
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var transactions []store.Transaction
	if err := json.Unmarshal([]byte(content), &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func extractPDFToText(ctx context.Context, filepath string) (string, error) {
	args := []string{
		"-layout",
		"-nopgbrk",
		filepath,
		"-",
	}
	cmd := exec.CommandContext(ctx, "pdftotext", args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
