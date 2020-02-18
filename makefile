.PHONY:test install

dir=.
testName=xxx
test:install clean
	go test -v $(dir) -run ^$(testName)$$

# 指定go版本，如：go1.14rc1。不过在使用前请确保本地已安装！
gov=go
gov_test:install clean
	$(gov) test -v $(dir) -run ^$(testName)$$

clean:
	go clean -testcache

install:
	go install ./...

# 通过git log查看文件的修改时间
file=xxx
file_last_ctime_from_git_log:
	git log -1 --pretty="format:%ci" $(file)
