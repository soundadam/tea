# frozen_string_literal: true

# The release workflow fills in @VERSION@ and the SHA-256 values and writes the result to
# Casks/tea.rb in soundadam/homebrew-tap.
cask "tea" do
  arch arm: "aarch64", intel: "x86_64"

  version "@VERSION@"
  sha256 arm:   "@MACOS_AARCH64_SHA256@",
         intel: "@MACOS_X86_64_SHA256@"

  url "https://github.com/soundadam/tea/releases/download/v#{version}/tea-#{version}-macos-#{arch}.tar.gz"
  name "tea"
  desc "Stay awake with the lid closed, then restore the previous sleep setting"
  homepage "https://soundadam.com/projects/tea/"

  livecheck do
    url :url
    strategy :github_latest
  end

  depends_on macos: :ventura

  binary "tea-#{version}-macos-#{arch}/tea"

  # Ad-hoc signed and not notarized: clear quarantine so Gatekeeper treats it like a self-updated copy.
  postflight_steps do
    run "/usr/bin/xattr", args: ["-dr", "com.apple.quarantine", "{{staged_path}}"]
  end

  caveats <<~EOS
    The passwordless helper is a copy of tea made by `tea auth register`.
    After upgrading, run `tea auth register` again (the client offers this
    as "Repair helper"); until then, privileged actions go through sudo.

    Before uninstalling, run `tea off` and `tea auth unregister`.
  EOS
end
