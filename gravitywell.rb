class Gravitywell < Formula
  desc "Turn a pool of Docker hosts into a single, virtual host"
  homepage "https://github.com/AlexsJones/gravitywell"
  url "https://github.com/AlexsJones/gravitywell/releases/download/v0.1.2/gravitywell_0.1.2_Darwin_amd64.tar.gz"
  sha256 "cb3d3c49cc0429e0f49d6425d661a41892a9513d32a191a0d5d05b77678f58fe"
  head "https://github.com/AlexsJones/gravitywell.git"

  def install
    bin.install "gravitywell"
  end

  test do
    system "#{bin}/gravitywell -h"
  end
end
