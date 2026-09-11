## Comment fd1aa53 gacha bolgan ishlarni va birinchi bussines logic ni techinical tushintirish

Birinchi main.go dan tushintirsam (/cmd/api) ni icida configlar yigiladi config faqat main.go asosiy filedan bajariladi. ikkinchi qadamda server  yaratiladi. va mux hamma routerlarni serverga beramiz. shundan keyin ularni run qilamiz.

main.go > database.new() -> Store.NewStore(db)(Bu yerda Product, Order/User repolari) -> services.NewServices()-> user.NewHandler(svc) -> app.ServeHTTP (bu yerda handler routerga boglanadi.)

Runtime tartibida, http sorov yuboriladi -> chi router (server/routes.go) -> middleware (Bu yerda logger, recoverer, RequireAuth va JWT tekshiruvi bor) -> Handler (jsoni oqish status code va error sodir bolsa uni togri error message bilan qaytarish kabi vazifalarni bajaradi.)-> Service ->biznes mantiq-> repository faqat SQL -> PosgreSQL ga tushadi 



## Comment 6eec0bf gacha bolgan ishlarni tushuntmasi

### CreateOrder /internal/order/repositry
CreateOrder Reposity chuntirishdan boshladim chunki handler va service kayerlar deyarli birhil, asosiy logika bu repository layerda, ishni boshlashda birinchi ishni tranzaksiyadan boshladim, chunki bitta functionda bir nechta database operationlar boladi jumladan, Order yaratish, Orderni Itemlarini yaratish, Stockni kamaytirish. endi detail bilan tushintiraman, birinchi tranzaksiya yaratilgandan song  QueryRowContext orqali bitta rowni olish uchun ishlatdim va  idempotency_keys mavjud bolsa bizga qayataradi shuningdek "user_id = $2" qismida oziga tegishli orderniyam olish logicasi yani idor protection ulangan. Undan keyin biz err sql.ErrNoRows shu qism order topilmadi degan qism shu yerda boladi.  shundan keyin Order yaratamiz INSERT qilish orqali.

 keyingi qadamda har bitta itemni insert qilib product id va quantyty qoshib insert qilamiz, va shundan keyin stockni kamaytimariz update va quantytni stockdaki quantitydan ayiramiz. RowsAffected orqali tekshiramiz row agar 0 bolsa ErrInsufficientStock hato tashaydi. va shundan keyin idempotency_keys yaratamiz. keyin tranzaksiyani commit qilib yopamiz


### GetOrder /internal/order/repositry
bu yerda logic odiyroq, faqat orderId va userid orqali idor qilish orqali userga tegishli orderni olamiz


### CancelOrder /internal/order/repositry

bu yerdaham bir nechta database operationlari bolganligi tufayli tranzaktion ochamiz. keyin Statusni tekshiramiz approved/rejected bols cancel qila olmaymiz., keyingi queryda stockdaki mahsulot sonini qaytarib qoyamiz. shudan keyingina order cancel database operation boladi. va tranzaktionni yopamiz commit qilish orqali