package shopify

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

func IsValid(id string) bool {
	for _, v := range Templates {
		if v == id {
			return true
		}
	}
	return false
}
