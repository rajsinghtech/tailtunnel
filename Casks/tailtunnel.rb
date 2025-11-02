cask "tailtunnel" do
  version "0.0.1"
  sha256 "PLACEHOLDER_SHA256"

  url "https://github.com/rajsinghtech/tailtunnel/releases/download/v#{version}/TailTunnel.zip"
  name "TailTunnel"
  desc "Web dashboard for Tailscale with browser-based SSH and network diagnostics"
  homepage "https://github.com/rajsinghtech/tailtunnel"

  livecheck do
    url :url
    strategy :github_latest
  end

  app "TailTunnel.app"

  postflight do
    system_command "/usr/bin/xattr",
                   args: ["-cr", "#{appdir}/TailTunnel.app"],
                   sudo: false
  end

  zap trash: "~/.tailtunnel"
end
