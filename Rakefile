task :default => [:push, :build]

GO = '~/Documents/go-cross-compile/go/bin/go'

def version
  return @version if @version
  
  open('gomilk.go') do |f|
    f.each do |line|
      if line =~ /const version = "(.*)"/
        @version = $1
        return @version
      end
    end
  end

  raise
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

desc "Push the repository"
task :push do
  system("git tag v#{version}")
  system("git push --tag")
  system("git push")
end

