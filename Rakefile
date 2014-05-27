task :default => :build

GO = '~/Documents/go-cross-compile/go/bin/go'

def version
  '0.2.0'
end

def zip_filename(os, fmt)
  "gomilk-#{version}-#{os}-#{fmt}.zip"
end

def build(os, fmt)
  filename = (os == 'windows') ? 'gomilk.exe' : 'gomilk'
  system("rm -f #{filename}")
  system("GOOS=#{os} GOARCH=#{fmt} #{GO} clean")
  system("GOOS=#{os} GOARCH=#{fmt} #{GO} build gomilk.go")
  system("zip archive/#{zip_filename(os, fmt)} #{filename}")
  puts "--> archive/#{zip_filename(os, fmt)}"
end

desc "Build the binaries"
task :build do
  build('darwin', 'amd64')
  build('windows', 'amd64')
  build('linux', 'amd64')
end

