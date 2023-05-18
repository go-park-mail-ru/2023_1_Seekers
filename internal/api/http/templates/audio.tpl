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
    <audio controls style="width: 100%;">
      <source src="{{.FilePath}}" width=>
      <p>
        Ваш браузер не поддерживает HTML5 аудио. Вы можете скачать данное аудио
      </p>
    </audio>
</main>
</body>
</html>