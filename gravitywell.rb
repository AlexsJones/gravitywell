class Gravitywell < Formula
  desc "Turn a pool of Docker hosts into a single, virtual host"
  homepage "https://github.com/AlexsJones/gravitywell"
  url "https://github.com/AlexsJones/gravitywell/archive/v0.1.2.tar.gz"
  sha256 "072580e0baf026b702847eb3cb7771c1630aee409b6b1d0ec8a36416784381c8"

  head "https://github.com/AlexsJones/gravitywell.git"

  bottle do
    cellar :any
    sha256 "6a8a15a511cd1bf5f10c9e53af8784d06d51fdd359ab533dfe8daaa7a6c872b3" => :yosemite
    sha256 "33f60f41343326f549c886ed60564a1672877f77c8636962d8c4c47c10d62f05" => :mavericks
    sha256 "e6b63640e4ff88d073c0449ea23a3c8c3ff46caffc72c3e4d4a2e159ffade10d" => :mountain_lion
  end

  depends_on "go" => :build

  def install
    mkdir_p buildpath/"src/github.com/AlexsJones"
    ln_s buildpath, buildpath/"src/github.com/AlexsJones/gravitywell"

    ENV["GOPATH"] = "#{buildpath}/Godeps/_workspace:#{buildpath}"

    Language::Go.stage_deps resources, buildpath/"src"

    system "go", "build", "-o", "gravitywell"

    bin.install "gravitywell"
  end

  test do
    output = shell_output(bin/"gravitywell --version")
    assert_match "swarm version #{version} (HEAD)", output
  end
end
