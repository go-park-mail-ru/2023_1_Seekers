<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>MailBx</title>
</head>
<body>
<main>
    <p><b>{{.FileName}} <a href="{{.FileDownload}}">скачать</a></b></p>
    <video src="{{.FilePath}}" controls autoplay id="my_video" style="width: 100vw; height: 90vh; object-fit: cover; position: fixed; top: 15; left: 0;">
      Ваш браузер не поддерживает просмотр видео,
      но Вы можете его скачать
    </video>
</main>
</body>
</html>
