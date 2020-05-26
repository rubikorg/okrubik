<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <title>Rubik Error</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      body {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        background-image: url('');
        background-image: url('https://i.imgur.com/tG8YsFb.png');
        background-repeat: repeat;
      }
      .gopher-rubik {
        position: absolute;
        width: 400px;
        bottom: 0;
        left: 0;
      }
      .content {
        margin-left: 25%;
        margin-right: 25%;
        margin-top: 2em;
        word-wrap: break-word;
      }
      /* .error-lines {
          padding: 0.5px;
      } */
      h2 {
          color: #AB256F;
      }

      footer {
          color: gray;
      }

      code {
          color: #25A6AB;
          font-weight: bold;
      }
    </style>
  </head>
  <body>
    <img class="gopher-rubik" src="https://i.imgur.com/Rx9VVDg.png" />
    <div class="content">
      <h2>{{ .Msg }}</h2>
      {{ range $element := .Stack }}
      <div class="error-lines">
        <p>❗️{{$element}}</p>
      </div>
      {{ end }}
      <footer>
          This view is only visible on default/development mode. It will no longer
          render if you set <code>RUBIK_ENV</code> other that empty or development.
      </footer>
    </div>
  </body>
</html>
