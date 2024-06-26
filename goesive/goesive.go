package goesive

import (
	"bytes"
	b64 "encoding/base64"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/goptos/goptos/io"
)

func check(e error) {
	if e != nil {
		log.Printf("%s", e)
		panic(e)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !(os.IsNotExist(err))
}

func copyFile(src string, dst string) {
	log.Printf("Copying: %q to %q", src, dst)
	index, err := io.ReadFile(src)
	check(err)
	err = io.WriteFile(dst, index)
	check(err)
}

func createFile(data string, filePath string) {
	log.Printf("Checking: %q", filePath)
	if exists(filePath) {
		log.Printf("Already exists: %q", filePath)
		return
	}

	log.Printf("Creating: %q", filePath)
	f, err := os.Create(filePath)
	check(err)
	defer f.Close()

	str, err := b64.StdEncoding.DecodeString(data)
	check(err)

	_, err = f.WriteString(string(str))
	check(err)
	f.Sync()
}

// func createPath(path string) {
// 	log.Printf("Creating: %q", path)
// 	err := io.WritePath(path)
// 	check(err)
// }

func goRun(dir string, args ...string) {
	var HOME = os.Getenv("HOME")
	var PATH = os.Getenv("PATH")
	var GOMODCACHE = os.Getenv("GOMODCACHE")
	var GOPATH = os.Getenv("GOPATH")
	var GOPTOS_VERBOSE = os.Getenv("GOPTOS_VERBOSE")
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	var cmd = exec.Command(args[0], args[1:]...)
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut
	cmd.Dir = dir
	cmd.Env = []string{
		"PATH=" + PATH,
		"HOME=" + HOME,
		"GOMODCACHE=" + GOMODCACHE,
		"GOPATH=" + GOPATH,
		"GOPTOS_VERBOSE=" + GOPTOS_VERBOSE,
		"GOOS=js",
		"GOARCH=wasm",
		"GONOPROXY=github.com/goptos"}
	err := cmd.Run()
	log.Printf("%s\n", strings.Join(cmd.Env[1:], " "))
	log.Printf("%s\n", strings.Join(cmd.Args, " "))
	if err != nil {
		log.Printf("%s", stdErr.String())
		log.Fatal(err)
	}
	if stdOut.String() != "" {
		log.Printf("%s", stdOut.String())
	}
}

func Build(src string) {
	goRun(src, "go", "generate", "-v")
	goRun(src, "go", "build", "-o", "../dist/main.wasm", "main.go")
}

func Pack(dist string) {
	var wasm = dist + "/main.wasm"
	log.Printf("Checking: %q", wasm)
	if !exists(wasm) {
		log.Fatalf("Cannot package until %q has first been built.", wasm)
	}
	log.Printf("Already exists: %q", wasm)
	createFile(dataIndexJs, dist+"/index.js")
	createFile(dataWasmExec, dist+"/wasm_exec.js")
	copyFile("index.html", dist+"/index.html")
}

func Serve(src string, dist string, port string) {
	var listen = "localhost:" + port
	Build(src)
	Pack(dist)
	log.Printf("Listening on http://%s...", listen)
	log.Fatal(http.ListenAndServe(listen, http.FileServer(http.Dir(dist))))
}

const dataIndexJs string = "Y29uc3QgZ28gPSBuZXcgR28oKTsKCldlYkFzc2VtYmx5Lmluc3RhbnRpYXRlU3RyZWFtaW5nKGZldGNoKCJtYWluLndhc20iKSwgZ28uaW1wb3J0T2JqZWN0KS50aGVuKGFzeW5jIChyZXN1bHQpID0+IHsKICBhd2FpdCBnby5ydW4ocmVzdWx0Lmluc3RhbmNlKTsKfSk7"

const dataWasmExec string = "Ly8gQ29weXJpZ2h0IDIwMTggVGhlIEdvIEF1dGhvcnMuIEFsbCByaWdodHMgcmVzZXJ2ZWQuCi8vIFVzZSBvZiB0aGlzIHNvdXJjZSBjb2RlIGlzIGdvdmVybmVkIGJ5IGEgQlNELXN0eWxlCi8vIGxpY2Vuc2UgdGhhdCBjYW4gYmUgZm91bmQgaW4gdGhlIExJQ0VOU0UgZmlsZS4KCi8vIGh0dHBzOi8vcmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbS9nb2xhbmcvZ28vbWFzdGVyL21pc2Mvd2FzbS93YXNtX2V4ZWMuanMKCiJ1c2Ugc3RyaWN0IjsKCigoKSA9PiB7Cgljb25zdCBlbm9zeXMgPSAoKSA9PiB7CgkJY29uc3QgZXJyID0gbmV3IEVycm9yKCJub3QgaW1wbGVtZW50ZWQiKTsKCQllcnIuY29kZSA9ICJFTk9TWVMiOwoJCXJldHVybiBlcnI7Cgl9OwoKCWlmICghZ2xvYmFsVGhpcy5mcykgewoJCWxldCBvdXRwdXRCdWYgPSAiIjsKCQlnbG9iYWxUaGlzLmZzID0gewoJCQljb25zdGFudHM6IHsgT19XUk9OTFk6IC0xLCBPX1JEV1I6IC0xLCBPX0NSRUFUOiAtMSwgT19UUlVOQzogLTEsIE9fQVBQRU5EOiAtMSwgT19FWENMOiAtMSB9LCAvLyB1bnVzZWQKCQkJd3JpdGVTeW5jKGZkLCBidWYpIHsKCQkJCW91dHB1dEJ1ZiArPSBkZWNvZGVyLmRlY29kZShidWYpOwoJCQkJY29uc3QgbmwgPSBvdXRwdXRCdWYubGFzdEluZGV4T2YoIlxuIik7CgkJCQlpZiAobmwgIT0gLTEpIHsKCQkJCQljb25zb2xlLmxvZyhvdXRwdXRCdWYuc3Vic3RyaW5nKDAsIG5sKSk7CgkJCQkJb3V0cHV0QnVmID0gb3V0cHV0QnVmLnN1YnN0cmluZyhubCArIDEpOwoJCQkJfQoJCQkJcmV0dXJuIGJ1Zi5sZW5ndGg7CgkJCX0sCgkJCXdyaXRlKGZkLCBidWYsIG9mZnNldCwgbGVuZ3RoLCBwb3NpdGlvbiwgY2FsbGJhY2spIHsKCQkJCWlmIChvZmZzZXQgIT09IDAgfHwgbGVuZ3RoICE9PSBidWYubGVuZ3RoIHx8IHBvc2l0aW9uICE9PSBudWxsKSB7CgkJCQkJY2FsbGJhY2soZW5vc3lzKCkpOwoJCQkJCXJldHVybjsKCQkJCX0KCQkJCWNvbnN0IG4gPSB0aGlzLndyaXRlU3luYyhmZCwgYnVmKTsKCQkJCWNhbGxiYWNrKG51bGwsIG4pOwoJCQl9LAoJCQljaG1vZChwYXRoLCBtb2RlLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCWNob3duKHBhdGgsIHVpZCwgZ2lkLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCWNsb3NlKGZkLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCWZjaG1vZChmZCwgbW9kZSwgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQlmY2hvd24oZmQsIHVpZCwgZ2lkLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCWZzdGF0KGZkLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCWZzeW5jKGZkLCBjYWxsYmFjaykgeyBjYWxsYmFjayhudWxsKTsgfSwKCQkJZnRydW5jYXRlKGZkLCBsZW5ndGgsIGNhbGxiYWNrKSB7IGNhbGxiYWNrKGVub3N5cygpKTsgfSwKCQkJbGNob3duKHBhdGgsIHVpZCwgZ2lkLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCWxpbmsocGF0aCwgbGluaywgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQlsc3RhdChwYXRoLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCW1rZGlyKHBhdGgsIHBlcm0sIGNhbGxiYWNrKSB7IGNhbGxiYWNrKGVub3N5cygpKTsgfSwKCQkJb3BlbihwYXRoLCBmbGFncywgbW9kZSwgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQlyZWFkKGZkLCBidWZmZXIsIG9mZnNldCwgbGVuZ3RoLCBwb3NpdGlvbiwgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQlyZWFkZGlyKHBhdGgsIGNhbGxiYWNrKSB7IGNhbGxiYWNrKGVub3N5cygpKTsgfSwKCQkJcmVhZGxpbmsocGF0aCwgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQlyZW5hbWUoZnJvbSwgdG8sIGNhbGxiYWNrKSB7IGNhbGxiYWNrKGVub3N5cygpKTsgfSwKCQkJcm1kaXIocGF0aCwgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQlzdGF0KHBhdGgsIGNhbGxiYWNrKSB7IGNhbGxiYWNrKGVub3N5cygpKTsgfSwKCQkJc3ltbGluayhwYXRoLCBsaW5rLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJCXRydW5jYXRlKHBhdGgsIGxlbmd0aCwgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQl1bmxpbmsocGF0aCwgY2FsbGJhY2spIHsgY2FsbGJhY2soZW5vc3lzKCkpOyB9LAoJCQl1dGltZXMocGF0aCwgYXRpbWUsIG10aW1lLCBjYWxsYmFjaykgeyBjYWxsYmFjayhlbm9zeXMoKSk7IH0sCgkJfTsKCX0KCglpZiAoIWdsb2JhbFRoaXMucHJvY2VzcykgewoJCWdsb2JhbFRoaXMucHJvY2VzcyA9IHsKCQkJZ2V0dWlkKCkgeyByZXR1cm4gLTE7IH0sCgkJCWdldGdpZCgpIHsgcmV0dXJuIC0xOyB9LAoJCQlnZXRldWlkKCkgeyByZXR1cm4gLTE7IH0sCgkJCWdldGVnaWQoKSB7IHJldHVybiAtMTsgfSwKCQkJZ2V0Z3JvdXBzKCkgeyB0aHJvdyBlbm9zeXMoKTsgfSwKCQkJcGlkOiAtMSwKCQkJcHBpZDogLTEsCgkJCXVtYXNrKCkgeyB0aHJvdyBlbm9zeXMoKTsgfSwKCQkJY3dkKCkgeyB0aHJvdyBlbm9zeXMoKTsgfSwKCQkJY2hkaXIoKSB7IHRocm93IGVub3N5cygpOyB9LAoJCX0KCX0KCglpZiAoIWdsb2JhbFRoaXMuY3J5cHRvKSB7CgkJdGhyb3cgbmV3IEVycm9yKCJnbG9iYWxUaGlzLmNyeXB0byBpcyBub3QgYXZhaWxhYmxlLCBwb2x5ZmlsbCByZXF1aXJlZCAoY3J5cHRvLmdldFJhbmRvbVZhbHVlcyBvbmx5KSIpOwoJfQoKCWlmICghZ2xvYmFsVGhpcy5wZXJmb3JtYW5jZSkgewoJCXRocm93IG5ldyBFcnJvcigiZ2xvYmFsVGhpcy5wZXJmb3JtYW5jZSBpcyBub3QgYXZhaWxhYmxlLCBwb2x5ZmlsbCByZXF1aXJlZCAocGVyZm9ybWFuY2Uubm93IG9ubHkpIik7Cgl9CgoJaWYgKCFnbG9iYWxUaGlzLlRleHRFbmNvZGVyKSB7CgkJdGhyb3cgbmV3IEVycm9yKCJnbG9iYWxUaGlzLlRleHRFbmNvZGVyIGlzIG5vdCBhdmFpbGFibGUsIHBvbHlmaWxsIHJlcXVpcmVkIik7Cgl9CgoJaWYgKCFnbG9iYWxUaGlzLlRleHREZWNvZGVyKSB7CgkJdGhyb3cgbmV3IEVycm9yKCJnbG9iYWxUaGlzLlRleHREZWNvZGVyIGlzIG5vdCBhdmFpbGFibGUsIHBvbHlmaWxsIHJlcXVpcmVkIik7Cgl9CgoJY29uc3QgZW5jb2RlciA9IG5ldyBUZXh0RW5jb2RlcigidXRmLTgiKTsKCWNvbnN0IGRlY29kZXIgPSBuZXcgVGV4dERlY29kZXIoInV0Zi04Iik7CgoJZ2xvYmFsVGhpcy5HbyA9IGNsYXNzIHsKCQljb25zdHJ1Y3RvcigpIHsKCQkJdGhpcy5hcmd2ID0gWyJqcyJdOwoJCQl0aGlzLmVudiA9IHt9OwoJCQl0aGlzLmV4aXQgPSAoY29kZSkgPT4gewoJCQkJaWYgKGNvZGUgIT09IDApIHsKCQkJCQljb25zb2xlLndhcm4oImV4aXQgY29kZToiLCBjb2RlKTsKCQkJCX0KCQkJfTsKCQkJdGhpcy5fZXhpdFByb21pc2UgPSBuZXcgUHJvbWlzZSgocmVzb2x2ZSkgPT4gewoJCQkJdGhpcy5fcmVzb2x2ZUV4aXRQcm9taXNlID0gcmVzb2x2ZTsKCQkJfSk7CgkJCXRoaXMuX3BlbmRpbmdFdmVudCA9IG51bGw7CgkJCXRoaXMuX3NjaGVkdWxlZFRpbWVvdXRzID0gbmV3IE1hcCgpOwoJCQl0aGlzLl9uZXh0Q2FsbGJhY2tUaW1lb3V0SUQgPSAxOwoKCQkJY29uc3Qgc2V0SW50NjQgPSAoYWRkciwgdikgPT4gewoJCQkJdGhpcy5tZW0uc2V0VWludDMyKGFkZHIgKyAwLCB2LCB0cnVlKTsKCQkJCXRoaXMubWVtLnNldFVpbnQzMihhZGRyICsgNCwgTWF0aC5mbG9vcih2IC8gNDI5NDk2NzI5NiksIHRydWUpOwoJCQl9CgoJCQljb25zdCBzZXRJbnQzMiA9IChhZGRyLCB2KSA9PiB7CgkJCQl0aGlzLm1lbS5zZXRVaW50MzIoYWRkciArIDAsIHYsIHRydWUpOwoJCQl9CgoJCQljb25zdCBnZXRJbnQ2NCA9IChhZGRyKSA9PiB7CgkJCQljb25zdCBsb3cgPSB0aGlzLm1lbS5nZXRVaW50MzIoYWRkciArIDAsIHRydWUpOwoJCQkJY29uc3QgaGlnaCA9IHRoaXMubWVtLmdldEludDMyKGFkZHIgKyA0LCB0cnVlKTsKCQkJCXJldHVybiBsb3cgKyBoaWdoICogNDI5NDk2NzI5NjsKCQkJfQoKCQkJY29uc3QgbG9hZFZhbHVlID0gKGFkZHIpID0+IHsKCQkJCWNvbnN0IGYgPSB0aGlzLm1lbS5nZXRGbG9hdDY0KGFkZHIsIHRydWUpOwoJCQkJaWYgKGYgPT09IDApIHsKCQkJCQlyZXR1cm4gdW5kZWZpbmVkOwoJCQkJfQoJCQkJaWYgKCFpc05hTihmKSkgewoJCQkJCXJldHVybiBmOwoJCQkJfQoKCQkJCWNvbnN0IGlkID0gdGhpcy5tZW0uZ2V0VWludDMyKGFkZHIsIHRydWUpOwoJCQkJcmV0dXJuIHRoaXMuX3ZhbHVlc1tpZF07CgkJCX0KCgkJCWNvbnN0IHN0b3JlVmFsdWUgPSAoYWRkciwgdikgPT4gewoJCQkJY29uc3QgbmFuSGVhZCA9IDB4N0ZGODAwMDA7CgoJCQkJaWYgKHR5cGVvZiB2ID09PSAibnVtYmVyIiAmJiB2ICE9PSAwKSB7CgkJCQkJaWYgKGlzTmFOKHYpKSB7CgkJCQkJCXRoaXMubWVtLnNldFVpbnQzMihhZGRyICsgNCwgbmFuSGVhZCwgdHJ1ZSk7CgkJCQkJCXRoaXMubWVtLnNldFVpbnQzMihhZGRyLCAwLCB0cnVlKTsKCQkJCQkJcmV0dXJuOwoJCQkJCX0KCQkJCQl0aGlzLm1lbS5zZXRGbG9hdDY0KGFkZHIsIHYsIHRydWUpOwoJCQkJCXJldHVybjsKCQkJCX0KCgkJCQlpZiAodiA9PT0gdW5kZWZpbmVkKSB7CgkJCQkJdGhpcy5tZW0uc2V0RmxvYXQ2NChhZGRyLCAwLCB0cnVlKTsKCQkJCQlyZXR1cm47CgkJCQl9CgoJCQkJbGV0IGlkID0gdGhpcy5faWRzLmdldCh2KTsKCQkJCWlmIChpZCA9PT0gdW5kZWZpbmVkKSB7CgkJCQkJaWQgPSB0aGlzLl9pZFBvb2wucG9wKCk7CgkJCQkJaWYgKGlkID09PSB1bmRlZmluZWQpIHsKCQkJCQkJaWQgPSB0aGlzLl92YWx1ZXMubGVuZ3RoOwoJCQkJCX0KCQkJCQl0aGlzLl92YWx1ZXNbaWRdID0gdjsKCQkJCQl0aGlzLl9nb1JlZkNvdW50c1tpZF0gPSAwOwoJCQkJCXRoaXMuX2lkcy5zZXQodiwgaWQpOwoJCQkJfQoJCQkJdGhpcy5fZ29SZWZDb3VudHNbaWRdKys7CgkJCQlsZXQgdHlwZUZsYWcgPSAwOwoJCQkJc3dpdGNoICh0eXBlb2YgdikgewoJCQkJCWNhc2UgIm9iamVjdCI6CgkJCQkJCWlmICh2ICE9PSBudWxsKSB7CgkJCQkJCQl0eXBlRmxhZyA9IDE7CgkJCQkJCX0KCQkJCQkJYnJlYWs7CgkJCQkJY2FzZSAic3RyaW5nIjoKCQkJCQkJdHlwZUZsYWcgPSAyOwoJCQkJCQlicmVhazsKCQkJCQljYXNlICJzeW1ib2wiOgoJCQkJCQl0eXBlRmxhZyA9IDM7CgkJCQkJCWJyZWFrOwoJCQkJCWNhc2UgImZ1bmN0aW9uIjoKCQkJCQkJdHlwZUZsYWcgPSA0OwoJCQkJCQlicmVhazsKCQkJCX0KCQkJCXRoaXMubWVtLnNldFVpbnQzMihhZGRyICsgNCwgbmFuSGVhZCB8IHR5cGVGbGFnLCB0cnVlKTsKCQkJCXRoaXMubWVtLnNldFVpbnQzMihhZGRyLCBpZCwgdHJ1ZSk7CgkJCX0KCgkJCWNvbnN0IGxvYWRTbGljZSA9IChhZGRyKSA9PiB7CgkJCQljb25zdCBhcnJheSA9IGdldEludDY0KGFkZHIgKyAwKTsKCQkJCWNvbnN0IGxlbiA9IGdldEludDY0KGFkZHIgKyA4KTsKCQkJCXJldHVybiBuZXcgVWludDhBcnJheSh0aGlzLl9pbnN0LmV4cG9ydHMubWVtLmJ1ZmZlciwgYXJyYXksIGxlbik7CgkJCX0KCgkJCWNvbnN0IGxvYWRTbGljZU9mVmFsdWVzID0gKGFkZHIpID0+IHsKCQkJCWNvbnN0IGFycmF5ID0gZ2V0SW50NjQoYWRkciArIDApOwoJCQkJY29uc3QgbGVuID0gZ2V0SW50NjQoYWRkciArIDgpOwoJCQkJY29uc3QgYSA9IG5ldyBBcnJheShsZW4pOwoJCQkJZm9yIChsZXQgaSA9IDA7IGkgPCBsZW47IGkrKykgewoJCQkJCWFbaV0gPSBsb2FkVmFsdWUoYXJyYXkgKyBpICogOCk7CgkJCQl9CgkJCQlyZXR1cm4gYTsKCQkJfQoKCQkJY29uc3QgbG9hZFN0cmluZyA9IChhZGRyKSA9PiB7CgkJCQljb25zdCBzYWRkciA9IGdldEludDY0KGFkZHIgKyAwKTsKCQkJCWNvbnN0IGxlbiA9IGdldEludDY0KGFkZHIgKyA4KTsKCQkJCXJldHVybiBkZWNvZGVyLmRlY29kZShuZXcgRGF0YVZpZXcodGhpcy5faW5zdC5leHBvcnRzLm1lbS5idWZmZXIsIHNhZGRyLCBsZW4pKTsKCQkJfQoKCQkJY29uc3QgdGltZU9yaWdpbiA9IERhdGUubm93KCkgLSBwZXJmb3JtYW5jZS5ub3coKTsKCQkJdGhpcy5pbXBvcnRPYmplY3QgPSB7CgkJCQlfZ290ZXN0OiB7CgkJCQkJYWRkOiAoYSwgYikgPT4gYSArIGIsCgkJCQl9LAoJCQkJZ29qczogewoJCQkJCS8vIEdvJ3MgU1AgZG9lcyBub3QgY2hhbmdlIGFzIGxvbmcgYXMgbm8gR28gY29kZSBpcyBydW5uaW5nLiBTb21lIG9wZXJhdGlvbnMgKGUuZy4gY2FsbHMsIGdldHRlcnMgYW5kIHNldHRlcnMpCgkJCQkJLy8gbWF5IHN5bmNocm9ub3VzbHkgdHJpZ2dlciBhIEdvIGV2ZW50IGhhbmRsZXIuIFRoaXMgbWFrZXMgR28gY29kZSBnZXQgZXhlY3V0ZWQgaW4gdGhlIG1pZGRsZSBvZiB0aGUgaW1wb3J0ZWQKCQkJCQkvLyBmdW5jdGlvbi4gQSBnb3JvdXRpbmUgY2FuIHN3aXRjaCB0byBhIG5ldyBzdGFjayBpZiB0aGUgY3VycmVudCBzdGFjayBpcyB0b28gc21hbGwgKHNlZSBtb3Jlc3RhY2sgZnVuY3Rpb24pLgoJCQkJCS8vIFRoaXMgY2hhbmdlcyB0aGUgU1AsIHRodXMgd2UgaGF2ZSB0byB1cGRhdGUgdGhlIFNQIHVzZWQgYnkgdGhlIGltcG9ydGVkIGZ1bmN0aW9uLgoKCQkJCQkvLyBmdW5jIHdhc21FeGl0KGNvZGUgaW50MzIpCgkJCQkJInJ1bnRpbWUud2FzbUV4aXQiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCBjb2RlID0gdGhpcy5tZW0uZ2V0SW50MzIoc3AgKyA4LCB0cnVlKTsKCQkJCQkJdGhpcy5leGl0ZWQgPSB0cnVlOwoJCQkJCQlkZWxldGUgdGhpcy5faW5zdDsKCQkJCQkJZGVsZXRlIHRoaXMuX3ZhbHVlczsKCQkJCQkJZGVsZXRlIHRoaXMuX2dvUmVmQ291bnRzOwoJCQkJCQlkZWxldGUgdGhpcy5faWRzOwoJCQkJCQlkZWxldGUgdGhpcy5faWRQb29sOwoJCQkJCQl0aGlzLmV4aXQoY29kZSk7CgkJCQkJfSwKCgkJCQkJLy8gZnVuYyB3YXNtV3JpdGUoZmQgdWludHB0ciwgcCB1bnNhZmUuUG9pbnRlciwgbiBpbnQzMikKCQkJCQkicnVudGltZS53YXNtV3JpdGUiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCBmZCA9IGdldEludDY0KHNwICsgOCk7CgkJCQkJCWNvbnN0IHAgPSBnZXRJbnQ2NChzcCArIDE2KTsKCQkJCQkJY29uc3QgbiA9IHRoaXMubWVtLmdldEludDMyKHNwICsgMjQsIHRydWUpOwoJCQkJCQlmcy53cml0ZVN5bmMoZmQsIG5ldyBVaW50OEFycmF5KHRoaXMuX2luc3QuZXhwb3J0cy5tZW0uYnVmZmVyLCBwLCBuKSk7CgkJCQkJfSwKCgkJCQkJLy8gZnVuYyByZXNldE1lbW9yeURhdGFWaWV3KCkKCQkJCQkicnVudGltZS5yZXNldE1lbW9yeURhdGFWaWV3IjogKHNwKSA9PiB7CgkJCQkJCXNwID4+Pj0gMDsKCQkJCQkJdGhpcy5tZW0gPSBuZXcgRGF0YVZpZXcodGhpcy5faW5zdC5leHBvcnRzLm1lbS5idWZmZXIpOwoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgbmFub3RpbWUxKCkgaW50NjQKCQkJCQkicnVudGltZS5uYW5vdGltZTEiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQlzZXRJbnQ2NChzcCArIDgsICh0aW1lT3JpZ2luICsgcGVyZm9ybWFuY2Uubm93KCkpICogMTAwMDAwMCk7CgkJCQkJfSwKCgkJCQkJLy8gZnVuYyB3YWxsdGltZSgpIChzZWMgaW50NjQsIG5zZWMgaW50MzIpCgkJCQkJInJ1bnRpbWUud2FsbHRpbWUiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCBtc2VjID0gKG5ldyBEYXRlKS5nZXRUaW1lKCk7CgkJCQkJCXNldEludDY0KHNwICsgOCwgbXNlYyAvIDEwMDApOwoJCQkJCQl0aGlzLm1lbS5zZXRJbnQzMihzcCArIDE2LCAobXNlYyAlIDEwMDApICogMTAwMDAwMCwgdHJ1ZSk7CgkJCQkJfSwKCgkJCQkJLy8gZnVuYyBzY2hlZHVsZVRpbWVvdXRFdmVudChkZWxheSBpbnQ2NCkgaW50MzIKCQkJCQkicnVudGltZS5zY2hlZHVsZVRpbWVvdXRFdmVudCI6IChzcCkgPT4gewoJCQkJCQlzcCA+Pj49IDA7CgkJCQkJCWNvbnN0IGlkID0gdGhpcy5fbmV4dENhbGxiYWNrVGltZW91dElEOwoJCQkJCQl0aGlzLl9uZXh0Q2FsbGJhY2tUaW1lb3V0SUQrKzsKCQkJCQkJdGhpcy5fc2NoZWR1bGVkVGltZW91dHMuc2V0KGlkLCBzZXRUaW1lb3V0KAoJCQkJCQkJKCkgPT4gewoJCQkJCQkJCXRoaXMuX3Jlc3VtZSgpOwoJCQkJCQkJCXdoaWxlICh0aGlzLl9zY2hlZHVsZWRUaW1lb3V0cy5oYXMoaWQpKSB7CgkJCQkJCQkJCS8vIGZvciBzb21lIHJlYXNvbiBHbyBmYWlsZWQgdG8gcmVnaXN0ZXIgdGhlIHRpbWVvdXQgZXZlbnQsIGxvZyBhbmQgdHJ5IGFnYWluCgkJCQkJCQkJCS8vICh0ZW1wb3Jhcnkgd29ya2Fyb3VuZCBmb3IgaHR0cHM6Ly9naXRodWIuY29tL2dvbGFuZy9nby9pc3N1ZXMvMjg5NzUpCgkJCQkJCQkJCWNvbnNvbGUud2Fybigic2NoZWR1bGVUaW1lb3V0RXZlbnQ6IG1pc3NlZCB0aW1lb3V0IGV2ZW50Iik7CgkJCQkJCQkJCXRoaXMuX3Jlc3VtZSgpOwoJCQkJCQkJCX0KCQkJCQkJCX0sCgkJCQkJCQlnZXRJbnQ2NChzcCArIDgpLAoJCQkJCQkpKTsKCQkJCQkJdGhpcy5tZW0uc2V0SW50MzIoc3AgKyAxNiwgaWQsIHRydWUpOwoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgY2xlYXJUaW1lb3V0RXZlbnQoaWQgaW50MzIpCgkJCQkJInJ1bnRpbWUuY2xlYXJUaW1lb3V0RXZlbnQiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCBpZCA9IHRoaXMubWVtLmdldEludDMyKHNwICsgOCwgdHJ1ZSk7CgkJCQkJCWNsZWFyVGltZW91dCh0aGlzLl9zY2hlZHVsZWRUaW1lb3V0cy5nZXQoaWQpKTsKCQkJCQkJdGhpcy5fc2NoZWR1bGVkVGltZW91dHMuZGVsZXRlKGlkKTsKCQkJCQl9LAoKCQkJCQkvLyBmdW5jIGdldFJhbmRvbURhdGEociBbXWJ5dGUpCgkJCQkJInJ1bnRpbWUuZ2V0UmFuZG9tRGF0YSI6IChzcCkgPT4gewoJCQkJCQlzcCA+Pj49IDA7CgkJCQkJCWNyeXB0by5nZXRSYW5kb21WYWx1ZXMobG9hZFNsaWNlKHNwICsgOCkpOwoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgZmluYWxpemVSZWYodiByZWYpCgkJCQkJInN5c2NhbGwvanMuZmluYWxpemVSZWYiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCBpZCA9IHRoaXMubWVtLmdldFVpbnQzMihzcCArIDgsIHRydWUpOwoJCQkJCQl0aGlzLl9nb1JlZkNvdW50c1tpZF0tLTsKCQkJCQkJaWYgKHRoaXMuX2dvUmVmQ291bnRzW2lkXSA9PT0gMCkgewoJCQkJCQkJY29uc3QgdiA9IHRoaXMuX3ZhbHVlc1tpZF07CgkJCQkJCQl0aGlzLl92YWx1ZXNbaWRdID0gbnVsbDsKCQkJCQkJCXRoaXMuX2lkcy5kZWxldGUodik7CgkJCQkJCQl0aGlzLl9pZFBvb2wucHVzaChpZCk7CgkJCQkJCX0KCQkJCQl9LAoKCQkJCQkvLyBmdW5jIHN0cmluZ1ZhbCh2YWx1ZSBzdHJpbmcpIHJlZgoJCQkJCSJzeXNjYWxsL2pzLnN0cmluZ1ZhbCI6IChzcCkgPT4gewoJCQkJCQlzcCA+Pj49IDA7CgkJCQkJCXN0b3JlVmFsdWUoc3AgKyAyNCwgbG9hZFN0cmluZyhzcCArIDgpKTsKCQkJCQl9LAoKCQkJCQkvLyBmdW5jIHZhbHVlR2V0KHYgcmVmLCBwIHN0cmluZykgcmVmCgkJCQkJInN5c2NhbGwvanMudmFsdWVHZXQiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCByZXN1bHQgPSBSZWZsZWN0LmdldChsb2FkVmFsdWUoc3AgKyA4KSwgbG9hZFN0cmluZyhzcCArIDE2KSk7CgkJCQkJCXNwID0gdGhpcy5faW5zdC5leHBvcnRzLmdldHNwKCkgPj4+IDA7IC8vIHNlZSBjb21tZW50IGFib3ZlCgkJCQkJCXN0b3JlVmFsdWUoc3AgKyAzMiwgcmVzdWx0KTsKCQkJCQl9LAoKCQkJCQkvLyBmdW5jIHZhbHVlU2V0KHYgcmVmLCBwIHN0cmluZywgeCByZWYpCgkJCQkJInN5c2NhbGwvanMudmFsdWVTZXQiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQlSZWZsZWN0LnNldChsb2FkVmFsdWUoc3AgKyA4KSwgbG9hZFN0cmluZyhzcCArIDE2KSwgbG9hZFZhbHVlKHNwICsgMzIpKTsKCQkJCQl9LAoKCQkJCQkvLyBmdW5jIHZhbHVlRGVsZXRlKHYgcmVmLCBwIHN0cmluZykKCQkJCQkic3lzY2FsbC9qcy52YWx1ZURlbGV0ZSI6IChzcCkgPT4gewoJCQkJCQlzcCA+Pj49IDA7CgkJCQkJCVJlZmxlY3QuZGVsZXRlUHJvcGVydHkobG9hZFZhbHVlKHNwICsgOCksIGxvYWRTdHJpbmcoc3AgKyAxNikpOwoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgdmFsdWVJbmRleCh2IHJlZiwgaSBpbnQpIHJlZgoJCQkJCSJzeXNjYWxsL2pzLnZhbHVlSW5kZXgiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQlzdG9yZVZhbHVlKHNwICsgMjQsIFJlZmxlY3QuZ2V0KGxvYWRWYWx1ZShzcCArIDgpLCBnZXRJbnQ2NChzcCArIDE2KSkpOwoJCQkJCX0sCgoJCQkJCS8vIHZhbHVlU2V0SW5kZXgodiByZWYsIGkgaW50LCB4IHJlZikKCQkJCQkic3lzY2FsbC9qcy52YWx1ZVNldEluZGV4IjogKHNwKSA9PiB7CgkJCQkJCXNwID4+Pj0gMDsKCQkJCQkJUmVmbGVjdC5zZXQobG9hZFZhbHVlKHNwICsgOCksIGdldEludDY0KHNwICsgMTYpLCBsb2FkVmFsdWUoc3AgKyAyNCkpOwoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgdmFsdWVDYWxsKHYgcmVmLCBtIHN0cmluZywgYXJncyBbXXJlZikgKHJlZiwgYm9vbCkKCQkJCQkic3lzY2FsbC9qcy52YWx1ZUNhbGwiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQl0cnkgewoJCQkJCQkJY29uc3QgdiA9IGxvYWRWYWx1ZShzcCArIDgpOwoJCQkJCQkJY29uc3QgbSA9IFJlZmxlY3QuZ2V0KHYsIGxvYWRTdHJpbmcoc3AgKyAxNikpOwoJCQkJCQkJY29uc3QgYXJncyA9IGxvYWRTbGljZU9mVmFsdWVzKHNwICsgMzIpOwoJCQkJCQkJY29uc3QgcmVzdWx0ID0gUmVmbGVjdC5hcHBseShtLCB2LCBhcmdzKTsKCQkJCQkJCXNwID0gdGhpcy5faW5zdC5leHBvcnRzLmdldHNwKCkgPj4+IDA7IC8vIHNlZSBjb21tZW50IGFib3ZlCgkJCQkJCQlzdG9yZVZhbHVlKHNwICsgNTYsIHJlc3VsdCk7CgkJCQkJCQl0aGlzLm1lbS5zZXRVaW50OChzcCArIDY0LCAxKTsKCQkJCQkJfSBjYXRjaCAoZXJyKSB7CgkJCQkJCQlzcCA9IHRoaXMuX2luc3QuZXhwb3J0cy5nZXRzcCgpID4+PiAwOyAvLyBzZWUgY29tbWVudCBhYm92ZQoJCQkJCQkJc3RvcmVWYWx1ZShzcCArIDU2LCBlcnIpOwoJCQkJCQkJdGhpcy5tZW0uc2V0VWludDgoc3AgKyA2NCwgMCk7CgkJCQkJCX0KCQkJCQl9LAoKCQkJCQkvLyBmdW5jIHZhbHVlSW52b2tlKHYgcmVmLCBhcmdzIFtdcmVmKSAocmVmLCBib29sKQoJCQkJCSJzeXNjYWxsL2pzLnZhbHVlSW52b2tlIjogKHNwKSA9PiB7CgkJCQkJCXNwID4+Pj0gMDsKCQkJCQkJdHJ5IHsKCQkJCQkJCWNvbnN0IHYgPSBsb2FkVmFsdWUoc3AgKyA4KTsKCQkJCQkJCWNvbnN0IGFyZ3MgPSBsb2FkU2xpY2VPZlZhbHVlcyhzcCArIDE2KTsKCQkJCQkJCWNvbnN0IHJlc3VsdCA9IFJlZmxlY3QuYXBwbHkodiwgdW5kZWZpbmVkLCBhcmdzKTsKCQkJCQkJCXNwID0gdGhpcy5faW5zdC5leHBvcnRzLmdldHNwKCkgPj4+IDA7IC8vIHNlZSBjb21tZW50IGFib3ZlCgkJCQkJCQlzdG9yZVZhbHVlKHNwICsgNDAsIHJlc3VsdCk7CgkJCQkJCQl0aGlzLm1lbS5zZXRVaW50OChzcCArIDQ4LCAxKTsKCQkJCQkJfSBjYXRjaCAoZXJyKSB7CgkJCQkJCQlzcCA9IHRoaXMuX2luc3QuZXhwb3J0cy5nZXRzcCgpID4+PiAwOyAvLyBzZWUgY29tbWVudCBhYm92ZQoJCQkJCQkJc3RvcmVWYWx1ZShzcCArIDQwLCBlcnIpOwoJCQkJCQkJdGhpcy5tZW0uc2V0VWludDgoc3AgKyA0OCwgMCk7CgkJCQkJCX0KCQkJCQl9LAoKCQkJCQkvLyBmdW5jIHZhbHVlTmV3KHYgcmVmLCBhcmdzIFtdcmVmKSAocmVmLCBib29sKQoJCQkJCSJzeXNjYWxsL2pzLnZhbHVlTmV3IjogKHNwKSA9PiB7CgkJCQkJCXNwID4+Pj0gMDsKCQkJCQkJdHJ5IHsKCQkJCQkJCWNvbnN0IHYgPSBsb2FkVmFsdWUoc3AgKyA4KTsKCQkJCQkJCWNvbnN0IGFyZ3MgPSBsb2FkU2xpY2VPZlZhbHVlcyhzcCArIDE2KTsKCQkJCQkJCWNvbnN0IHJlc3VsdCA9IFJlZmxlY3QuY29uc3RydWN0KHYsIGFyZ3MpOwoJCQkJCQkJc3AgPSB0aGlzLl9pbnN0LmV4cG9ydHMuZ2V0c3AoKSA+Pj4gMDsgLy8gc2VlIGNvbW1lbnQgYWJvdmUKCQkJCQkJCXN0b3JlVmFsdWUoc3AgKyA0MCwgcmVzdWx0KTsKCQkJCQkJCXRoaXMubWVtLnNldFVpbnQ4KHNwICsgNDgsIDEpOwoJCQkJCQl9IGNhdGNoIChlcnIpIHsKCQkJCQkJCXNwID0gdGhpcy5faW5zdC5leHBvcnRzLmdldHNwKCkgPj4+IDA7IC8vIHNlZSBjb21tZW50IGFib3ZlCgkJCQkJCQlzdG9yZVZhbHVlKHNwICsgNDAsIGVycik7CgkJCQkJCQl0aGlzLm1lbS5zZXRVaW50OChzcCArIDQ4LCAwKTsKCQkJCQkJfQoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgdmFsdWVMZW5ndGgodiByZWYpIGludAoJCQkJCSJzeXNjYWxsL2pzLnZhbHVlTGVuZ3RoIjogKHNwKSA9PiB7CgkJCQkJCXNwID4+Pj0gMDsKCQkJCQkJc2V0SW50NjQoc3AgKyAxNiwgcGFyc2VJbnQobG9hZFZhbHVlKHNwICsgOCkubGVuZ3RoKSk7CgkJCQkJfSwKCgkJCQkJLy8gdmFsdWVQcmVwYXJlU3RyaW5nKHYgcmVmKSAocmVmLCBpbnQpCgkJCQkJInN5c2NhbGwvanMudmFsdWVQcmVwYXJlU3RyaW5nIjogKHNwKSA9PiB7CgkJCQkJCXNwID4+Pj0gMDsKCQkJCQkJY29uc3Qgc3RyID0gZW5jb2Rlci5lbmNvZGUoU3RyaW5nKGxvYWRWYWx1ZShzcCArIDgpKSk7CgkJCQkJCXN0b3JlVmFsdWUoc3AgKyAxNiwgc3RyKTsKCQkJCQkJc2V0SW50NjQoc3AgKyAyNCwgc3RyLmxlbmd0aCk7CgkJCQkJfSwKCgkJCQkJLy8gdmFsdWVMb2FkU3RyaW5nKHYgcmVmLCBiIFtdYnl0ZSkKCQkJCQkic3lzY2FsbC9qcy52YWx1ZUxvYWRTdHJpbmciOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCBzdHIgPSBsb2FkVmFsdWUoc3AgKyA4KTsKCQkJCQkJbG9hZFNsaWNlKHNwICsgMTYpLnNldChzdHIpOwoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgdmFsdWVJbnN0YW5jZU9mKHYgcmVmLCB0IHJlZikgYm9vbAoJCQkJCSJzeXNjYWxsL2pzLnZhbHVlSW5zdGFuY2VPZiI6IChzcCkgPT4gewoJCQkJCQlzcCA+Pj49IDA7CgkJCQkJCXRoaXMubWVtLnNldFVpbnQ4KHNwICsgMjQsIChsb2FkVmFsdWUoc3AgKyA4KSBpbnN0YW5jZW9mIGxvYWRWYWx1ZShzcCArIDE2KSkgPyAxIDogMCk7CgkJCQkJfSwKCgkJCQkJLy8gZnVuYyBjb3B5Qnl0ZXNUb0dvKGRzdCBbXWJ5dGUsIHNyYyByZWYpIChpbnQsIGJvb2wpCgkJCQkJInN5c2NhbGwvanMuY29weUJ5dGVzVG9HbyI6IChzcCkgPT4gewoJCQkJCQlzcCA+Pj49IDA7CgkJCQkJCWNvbnN0IGRzdCA9IGxvYWRTbGljZShzcCArIDgpOwoJCQkJCQljb25zdCBzcmMgPSBsb2FkVmFsdWUoc3AgKyAzMik7CgkJCQkJCWlmICghKHNyYyBpbnN0YW5jZW9mIFVpbnQ4QXJyYXkgfHwgc3JjIGluc3RhbmNlb2YgVWludDhDbGFtcGVkQXJyYXkpKSB7CgkJCQkJCQl0aGlzLm1lbS5zZXRVaW50OChzcCArIDQ4LCAwKTsKCQkJCQkJCXJldHVybjsKCQkJCQkJfQoJCQkJCQljb25zdCB0b0NvcHkgPSBzcmMuc3ViYXJyYXkoMCwgZHN0Lmxlbmd0aCk7CgkJCQkJCWRzdC5zZXQodG9Db3B5KTsKCQkJCQkJc2V0SW50NjQoc3AgKyA0MCwgdG9Db3B5Lmxlbmd0aCk7CgkJCQkJCXRoaXMubWVtLnNldFVpbnQ4KHNwICsgNDgsIDEpOwoJCQkJCX0sCgoJCQkJCS8vIGZ1bmMgY29weUJ5dGVzVG9KUyhkc3QgcmVmLCBzcmMgW11ieXRlKSAoaW50LCBib29sKQoJCQkJCSJzeXNjYWxsL2pzLmNvcHlCeXRlc1RvSlMiOiAoc3ApID0+IHsKCQkJCQkJc3AgPj4+PSAwOwoJCQkJCQljb25zdCBkc3QgPSBsb2FkVmFsdWUoc3AgKyA4KTsKCQkJCQkJY29uc3Qgc3JjID0gbG9hZFNsaWNlKHNwICsgMTYpOwoJCQkJCQlpZiAoIShkc3QgaW5zdGFuY2VvZiBVaW50OEFycmF5IHx8IGRzdCBpbnN0YW5jZW9mIFVpbnQ4Q2xhbXBlZEFycmF5KSkgewoJCQkJCQkJdGhpcy5tZW0uc2V0VWludDgoc3AgKyA0OCwgMCk7CgkJCQkJCQlyZXR1cm47CgkJCQkJCX0KCQkJCQkJY29uc3QgdG9Db3B5ID0gc3JjLnN1YmFycmF5KDAsIGRzdC5sZW5ndGgpOwoJCQkJCQlkc3Quc2V0KHRvQ29weSk7CgkJCQkJCXNldEludDY0KHNwICsgNDAsIHRvQ29weS5sZW5ndGgpOwoJCQkJCQl0aGlzLm1lbS5zZXRVaW50OChzcCArIDQ4LCAxKTsKCQkJCQl9LAoKCQkJCQkiZGVidWciOiAodmFsdWUpID0+IHsKCQkJCQkJY29uc29sZS5sb2codmFsdWUpOwoJCQkJCX0sCgkJCQl9CgkJCX07CgkJfQoKCQlhc3luYyBydW4oaW5zdGFuY2UpIHsKCQkJaWYgKCEoaW5zdGFuY2UgaW5zdGFuY2VvZiBXZWJBc3NlbWJseS5JbnN0YW5jZSkpIHsKCQkJCXRocm93IG5ldyBFcnJvcigiR28ucnVuOiBXZWJBc3NlbWJseS5JbnN0YW5jZSBleHBlY3RlZCIpOwoJCQl9CgkJCXRoaXMuX2luc3QgPSBpbnN0YW5jZTsKCQkJdGhpcy5tZW0gPSBuZXcgRGF0YVZpZXcodGhpcy5faW5zdC5leHBvcnRzLm1lbS5idWZmZXIpOwoJCQl0aGlzLl92YWx1ZXMgPSBbIC8vIEpTIHZhbHVlcyB0aGF0IEdvIGN1cnJlbnRseSBoYXMgcmVmZXJlbmNlcyB0bywgaW5kZXhlZCBieSByZWZlcmVuY2UgaWQKCQkJCU5hTiwKCQkJCTAsCgkJCQludWxsLAoJCQkJdHJ1ZSwKCQkJCWZhbHNlLAoJCQkJZ2xvYmFsVGhpcywKCQkJCXRoaXMsCgkJCV07CgkJCXRoaXMuX2dvUmVmQ291bnRzID0gbmV3IEFycmF5KHRoaXMuX3ZhbHVlcy5sZW5ndGgpLmZpbGwoSW5maW5pdHkpOyAvLyBudW1iZXIgb2YgcmVmZXJlbmNlcyB0aGF0IEdvIGhhcyB0byBhIEpTIHZhbHVlLCBpbmRleGVkIGJ5IHJlZmVyZW5jZSBpZAoJCQl0aGlzLl9pZHMgPSBuZXcgTWFwKFsgLy8gbWFwcGluZyBmcm9tIEpTIHZhbHVlcyB0byByZWZlcmVuY2UgaWRzCgkJCQlbMCwgMV0sCgkJCQlbbnVsbCwgMl0sCgkJCQlbdHJ1ZSwgM10sCgkJCQlbZmFsc2UsIDRdLAoJCQkJW2dsb2JhbFRoaXMsIDVdLAoJCQkJW3RoaXMsIDZdLAoJCQldKTsKCQkJdGhpcy5faWRQb29sID0gW107ICAgLy8gdW51c2VkIGlkcyB0aGF0IGhhdmUgYmVlbiBnYXJiYWdlIGNvbGxlY3RlZAoJCQl0aGlzLmV4aXRlZCA9IGZhbHNlOyAvLyB3aGV0aGVyIHRoZSBHbyBwcm9ncmFtIGhhcyBleGl0ZWQKCgkJCS8vIFBhc3MgY29tbWFuZCBsaW5lIGFyZ3VtZW50cyBhbmQgZW52aXJvbm1lbnQgdmFyaWFibGVzIHRvIFdlYkFzc2VtYmx5IGJ5IHdyaXRpbmcgdGhlbSB0byB0aGUgbGluZWFyIG1lbW9yeS4KCQkJbGV0IG9mZnNldCA9IDQwOTY7CgoJCQljb25zdCBzdHJQdHIgPSAoc3RyKSA9PiB7CgkJCQljb25zdCBwdHIgPSBvZmZzZXQ7CgkJCQljb25zdCBieXRlcyA9IGVuY29kZXIuZW5jb2RlKHN0ciArICJcMCIpOwoJCQkJbmV3IFVpbnQ4QXJyYXkodGhpcy5tZW0uYnVmZmVyLCBvZmZzZXQsIGJ5dGVzLmxlbmd0aCkuc2V0KGJ5dGVzKTsKCQkJCW9mZnNldCArPSBieXRlcy5sZW5ndGg7CgkJCQlpZiAob2Zmc2V0ICUgOCAhPT0gMCkgewoJCQkJCW9mZnNldCArPSA4IC0gKG9mZnNldCAlIDgpOwoJCQkJfQoJCQkJcmV0dXJuIHB0cjsKCQkJfTsKCgkJCWNvbnN0IGFyZ2MgPSB0aGlzLmFyZ3YubGVuZ3RoOwoKCQkJY29uc3QgYXJndlB0cnMgPSBbXTsKCQkJdGhpcy5hcmd2LmZvckVhY2goKGFyZykgPT4gewoJCQkJYXJndlB0cnMucHVzaChzdHJQdHIoYXJnKSk7CgkJCX0pOwoJCQlhcmd2UHRycy5wdXNoKDApOwoKCQkJY29uc3Qga2V5cyA9IE9iamVjdC5rZXlzKHRoaXMuZW52KS5zb3J0KCk7CgkJCWtleXMuZm9yRWFjaCgoa2V5KSA9PiB7CgkJCQlhcmd2UHRycy5wdXNoKHN0clB0cihgJHtrZXl9PSR7dGhpcy5lbnZba2V5XX1gKSk7CgkJCX0pOwoJCQlhcmd2UHRycy5wdXNoKDApOwoKCQkJY29uc3QgYXJndiA9IG9mZnNldDsKCQkJYXJndlB0cnMuZm9yRWFjaCgocHRyKSA9PiB7CgkJCQl0aGlzLm1lbS5zZXRVaW50MzIob2Zmc2V0LCBwdHIsIHRydWUpOwoJCQkJdGhpcy5tZW0uc2V0VWludDMyKG9mZnNldCArIDQsIDAsIHRydWUpOwoJCQkJb2Zmc2V0ICs9IDg7CgkJCX0pOwoKCQkJLy8gVGhlIGxpbmtlciBndWFyYW50ZWVzIGdsb2JhbCBkYXRhIHN0YXJ0cyBmcm9tIGF0IGxlYXN0IHdhc21NaW5EYXRhQWRkci4KCQkJLy8gS2VlcCBpbiBzeW5jIHdpdGggY21kL2xpbmsvaW50ZXJuYWwvbGQvZGF0YS5nbzp3YXNtTWluRGF0YUFkZHIuCgkJCWNvbnN0IHdhc21NaW5EYXRhQWRkciA9IDQwOTYgKyA4MTkyOwoJCQlpZiAob2Zmc2V0ID49IHdhc21NaW5EYXRhQWRkcikgewoJCQkJdGhyb3cgbmV3IEVycm9yKCJ0b3RhbCBsZW5ndGggb2YgY29tbWFuZCBsaW5lIGFuZCBlbnZpcm9ubWVudCB2YXJpYWJsZXMgZXhjZWVkcyBsaW1pdCIpOwoJCQl9CgoJCQl0aGlzLl9pbnN0LmV4cG9ydHMucnVuKGFyZ2MsIGFyZ3YpOwoJCQlpZiAodGhpcy5leGl0ZWQpIHsKCQkJCXRoaXMuX3Jlc29sdmVFeGl0UHJvbWlzZSgpOwoJCQl9CgkJCWF3YWl0IHRoaXMuX2V4aXRQcm9taXNlOwoJCX0KCgkJX3Jlc3VtZSgpIHsKCQkJaWYgKHRoaXMuZXhpdGVkKSB7CgkJCQl0aHJvdyBuZXcgRXJyb3IoIkdvIHByb2dyYW0gaGFzIGFscmVhZHkgZXhpdGVkIik7CgkJCX0KCQkJdGhpcy5faW5zdC5leHBvcnRzLnJlc3VtZSgpOwoJCQlpZiAodGhpcy5leGl0ZWQpIHsKCQkJCXRoaXMuX3Jlc29sdmVFeGl0UHJvbWlzZSgpOwoJCQl9CgkJfQoKCQlfbWFrZUZ1bmNXcmFwcGVyKGlkKSB7CgkJCWNvbnN0IGdvID0gdGhpczsKCQkJcmV0dXJuIGZ1bmN0aW9uICgpIHsKCQkJCWNvbnN0IGV2ZW50ID0geyBpZDogaWQsIHRoaXM6IHRoaXMsIGFyZ3M6IGFyZ3VtZW50cyB9OwoJCQkJZ28uX3BlbmRpbmdFdmVudCA9IGV2ZW50OwoJCQkJZ28uX3Jlc3VtZSgpOwoJCQkJcmV0dXJuIGV2ZW50LnJlc3VsdDsKCQkJfTsKCQl9Cgl9Cn0pKCk7"
