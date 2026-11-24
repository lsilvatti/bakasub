# Getting an OpenSubtitles API key

When a video has no subtitle track to translate, BakaSub can search [OpenSubtitles](https://www.opensubtitles.com/) for one and translate from it. This is optional: you can also pick a subtitle file you already have.

Each person uses their own API key, and it is free.

## Steps

1. Create an account at [opensubtitles.com](https://www.opensubtitles.com/) and confirm your email.
2. Open the **API consumers** page at [opensubtitles.com/consumers](https://www.opensubtitles.com/consumers).
3. Create a new consumer. A name is enough, for example `BakaSub`.
4. Copy the API key shown for that consumer.
5. In BakaSub, paste it into **Settings → Keys → OpenSubtitles**. BakaSub checks the key with OpenSubtitles before saving.

OpenSubtitles may rename buttons or pages over time. The consumers page is always inside your account.

## Download limit and your account

OpenSubtitles limits how many subtitles you can download per day. Searching is not counted, only downloading. BakaSub keeps what it downloaded, so picking the same subtitle again costs nothing.

If you want a bigger limit, sign in with your OpenSubtitles account in the same Settings section. The limit becomes the one of your account. Use your username, because OpenSubtitles does not accept the email at sign-in.

## Which result to pick

- Prefer the results marked "Made for this file". That subtitle was made for this exact video file, so the timing matches.
- The other results match by title only and may be out of sync with your video. BakaSub does not adjust timing.
- Avoid the ones marked "Machine translated". Translating a translation loses meaning twice.

## Privacy

- The API key and the password are stored encrypted on your computer, in your system's keychain when there is one. They leave it only in requests to OpenSubtitles.
- A search sends the title, season and episode taken from the file name. It also sends a fingerprint of the video file, computed from its size and a small part of its content. The video itself is never uploaded.
