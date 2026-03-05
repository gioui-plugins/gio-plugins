include ../../makefile.mk

js:
	$(GO) run $(GOGIO) -target js -o ./wasm .

android:
	PATH=$(ANDROID_JAVA_ROOT):$(PATH) ANDROID_SDK_ROOT=$(ANDROID_SDK_ROOT) $(GO) run $(GOGIO) -minsdk 23 -x -appid "com.gio.test" -target android -o ./android.apk .

ios:
	$(GO) run $(GOGIO) -work -minsdk 15 -target ios -appid "gio.demo.inkeliz.com" -signkey "../../iosgiodemopublish.mobileprovision" -o ./ios.ipa .

windows:
	$(GO) run $(GOGIO) -target windows -arch amd64 -o ./windows_amd64.exe .

macos:
	$(GO) run $(GOGIO) -target macos -arch arm64 -appid "gio.demo.inkeliz.com" -x -target macos -signkey "../../macosgiodemoprofile.provisionprofile" -o ./macos_arm64.app .

all: js android ios windows macos