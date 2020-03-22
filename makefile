.PHONY:test install build gen

# 指定go版本，如：go1.14rc1。不过在使用前请确保本地已安装！
gov=go

# 开启调试
debug=

dir=.
testName=Test*
test:install clean
	env debug=$(debug) $(gov) test -v $(dir) -run ^$(testName)$$

clean:
	$(gov) clean -testcache

install:
	$(gov) install ./...

build:
	cd cmd/gen && \
	$(gov) install

gen:install build
	gen -r

# 通过git log查看文件的修改时间
file=xxx
file_last_ctime_from_git_log:
	git log -1 --pretty="format:%ci" $(file)
