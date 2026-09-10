Comment fd1aa53 gacha bolgan ishlarni va birinchi bussines logic ni techinical tushintirish

Birinchi main.go dan tushintirsam (/cmd/api) ni icida configlar yigiladi config faqat main.go asosiy filedan bajariladi. ikkinchi qadamda server  yaratiladi. va mux hamma routerlarni serverga beramiz. shundan keyin ularni run qilamiz.

main.go > database.new() -> Store.NewStore(db)(Bu yerda Product, Order/User repolari) -> services.NewServices()-> user.NewHandler(svc) -> app.ServeHTTP (bu yerda handler routerga boglanadi.)

Runtime tartibida, http sorov yuboriladi -> chi router (server/routes.go) -> middleware (Bu yerda logger, recoverer, RequireAuth va JWT tekshiruvi bor) -> Handler (jsoni oqish status code va error sodir bolsa uni togri error message bilan qaytarish kabi vazifalarni bajaradi.)-> Service ->biznes mantiq-> repository faqat SQL -> PosgreSQL ga tushadi 
