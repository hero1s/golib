
func main() {
	uIntOpt := grpc.UnaryInterceptor(panichandler.UnaryPanicHandler)
	sIntOpt := grpc.StreamInterceptor(panichandler.StreamPanicHandler)
	grpc.NewServer(uIntOpt, sIntOpt)
}
func main() {
	panichandler.InstallPanicHandler(func(r interface{}) {
		fmt.Printf("panic happened: %v", r)
	}
}