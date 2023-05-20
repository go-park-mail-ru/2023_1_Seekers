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
    <p><b style="text-align: center">{{.FileName}}   <a href="{{.FileDownload}}">скачать</a></b></p>
    <iframe src='https://docs.google.com/viewer?embedded=true&url={{.FilePath}}' width='100%' height='800px' frameborder='0'></iframe>
</main>
</body>
</html>
