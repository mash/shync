package shopify

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/chromedp/chromedp"
)

var (
	Templates = []string{
		"order_confirmation",              // 注文の確認
		"order_edited",                    // 注文の編集
		"order_edit_invoice",              // 注文編集済みの請求書
		"order_invoice",                   // 注文の請求書
		"order_cancelled",                 // 注文のキャンセル
		"refund_notification",             // 注文の返金
		"draft_order_invoice",             // 下書き注文の請求書
		"buy_online",                      // POSからのメールカート
		"abandoned_checkout_notification", // カゴ落ち
		"pos_exchange_receipt",            // POS交換レシート
		"gift_card_notification",          // ギフトカードの作成
		"failed_payment_processing",       // 支払いエラー
		"fulfillment_request",             // フルフィルメントのリクエスト
		"shipping_confirmation",           // 配送情報通知
		"shipping_update",                 // 配送更新
		"shipment_out_for_delivery",       // 配達中
		"shipment_delivered",              // 配達済み
		"local_out_for_delivery",          // 配達中
		"local_delivered",                 // 配達済み
		"local_missed_delivery",           // 不在配達
		"ready_for_pickup",                // 受取の準備完了
		"pickup_receipt",                  // 店頭受取済み
		"customer_account_activate",       // お客様アカウントの招待
		"customer_account_welcome",        // お客様アカウントへの挨拶
		"customer_account_reset",          // お客様アカウントのパスワードのリセット
		"customer_update_payment_method",  // お客様による決済方法更新のリクエスト
		"contact_buyer",                   // お客様への連絡
		"customer_marketing_confirmation", // 確認メール
		"return_created",                  // 返品の手順
		"return_label_notification",       // 返品用ラベルの手順
		"new_order_notification",          // 新しい注文
	}
)

// head=false means headless
func ChromeContext(head bool) (context.Context, func()) {
	var ctx context.Context
	var cancel func()
	if !head {
		// headless
		ctx, cancel = chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	} else {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false), // headless=false に変更
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-extensions", false),
			chromedp.Flag("hide-scrollbars", false),
			chromedp.Flag("mute-audio", false),
		)

		allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)

		ctx_, cancel2 := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf)) // cancel() を呼ばないように変更
		ctx = ctx_
		cancel = func() {
			cancel1()
			cancel2()
		}
	}
	return ctx, cancel
}

func Login(shop, username, password string, location *string) (chromedp.Tasks, error) {
	submit := `//button[@type='submit']`
	firstStore := `//input[@name='shop[domain]']`
	secondEmail := `//input[@name='account[email]']`
	thirdPassword := `//input[@name='account[password]']`
	linkSettings := `//a[@href='/admin/settings']`

	return chromedp.Tasks{
		chromedp.Navigate(`https://accounts.shopify.com/store-login?new_store_login=true`),
		chromedp.WaitVisible(firstStore),
		chromedp.SendKeys(firstStore, shop),
		chromedp.Submit(submit),

		chromedp.WaitVisible(secondEmail),
		chromedp.SendKeys(secondEmail, username),
		chromedp.Submit(submit),

		chromedp.WaitVisible(thirdPassword),
		chromedp.SendKeys(thirdPassword, password),
		chromedp.Submit(submit),

		chromedp.WaitVisible(linkSettings),
		chromedp.Location(location),
	}, nil
}

// id must be one of Templates
// The email template body will be set to *body
func FetchEmailTemplate(location, id string, subject, body *string) (chromedp.Tasks, error) {
	found := false
	for _, v := range Templates {
		if v == id {
			found = true
		}
	}
	if !found {
		// program error is fatal
		return nil, fmt.Errorf("FetchEmailTemplate: invalid email template id: %s", id)
	}

	path, err := url.Parse(fmt.Sprintf("/admin/email_templates/%s/edit", id))
	if err != nil {
		return nil, fmt.Errorf("FetchEmailTemplate: %w", err)
	}
	base, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("FetchEmailTemplate: %w", err)
	}
	next := base.ResolveReference(path)

	input := `//input[@name='email_template[title]']`
	textarea := `//textarea[@name='email_template[body_html]']`
	return chromedp.Tasks{
		chromedp.Navigate(next.String()),
		chromedp.WaitVisible(textarea),
		chromedp.Value(input, subject),
		chromedp.Value(textarea, body),
	}, nil
}
