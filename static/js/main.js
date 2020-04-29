(function ($) {
  var Game = function () {
    this.GameData = {};
    this.GameAlreadyInProgressBool = false;
    this.PlayerIndex = -1;

    $("#play-area").hide();


    this.Render = function(data) {
      $("#debug").text(data);
    },
    this.PollGame = function() {
      $this = this;
      window.setInterval(function () {
        $.getJSON( "/api/game", function( data ) {
          $this.GameData = data;
          $("#score tbody").empty()
          for (var i in $this.GameData.Players) {
            player = $this.GameData.Players[i];
    
            scores = $('<tr><td>' + player.Name.replace("[", "").replace("]", "") + '</td><td>' + player.Score + '</td></tr>');
            $("#score tbody").append(scores);
          }

          if (localStorage.getItem("game-id") != data.Id) {
            localStorage.setItem("game-id", data.Id);
            location.reload();
          }

          if (data.Started) {
            $this.GameAlreadyInProgress();
          } else {
            $("#lobby").fadeIn(1000);
          }
          
        });
      }, 1000)
    },
    this.New = function () {
      $this = this;
      $.getJSON( "/api/game/new", function( data ) {
        $this.GameData = data;
        localStorage.setItem("game-id", data.Id);
      });
      localStorage.removeItem("player-id");
      
    },
    this.AddPlayer = function (name) {
      $this = this;
      $.getJSON( "/api/game/players/new?name=" + name, function( data ) {
        localStorage.setItem("player-id", data.Id);
      });
    },
    this.Start = function () {
      $this = this;
      $.getJSON( "/api/game/start", function( data ) {
        $this.GameData = data;
        $this.RenderCards();
        $("#play-area").fadeIn();
      });
    },
    this.RenderCards = function () {
      $this = this;

      darkCardText = $this.GameData.BlackCard.Text
      darkCard = $('<div class="card dark"><p class="text">' + darkCardText + '</p></div>');
      $("#dark-card").empty().append(darkCard);


      for (var i in $this.GameData.Players) {
        player = $this.GameData.Players[i];
        if (player.Id == localStorage.getItem("player-id")) {
          $this.PlayerIndex = i;
        }
      } 
    
      if (this.GameData.Players[$this.PlayerIndex].IsCzar) {
        text = "You have been selected as the Card Czar.";
        $("#czar").empty();
        $("#czar").append($('<div class="czar-text">' + text + '</div>')).show();

        cardsInPlayCount = -1

        window.setInterval(function () {
            $("#cards").empty();
            for (var i in $this.GameData.CardsInPlay) {
              text = $this.GameData.CardsInPlay[i].Card.Text;
              id = $this.GameData.CardsInPlay[i].Card.Id;
              card = $('<div class="card light"><p class="text">' + text + '</p><a href="#" data-payload="?card=' + id + '" class="award">Award</a></div>');
              $("#cards").append(card);
            }
        }, 2000);
      } else {
        for (var i in this.GameData.Players[$this.PlayerIndex].Cards) {
          text = $this.GameData.Players[$this.PlayerIndex].Cards[i].Text;
          id = $this.GameData.Players[$this.PlayerIndex].Cards[i].Id;

          if ($this.GameData.Players[$this.PlayerIndex].Cards[i].Played) {
            card = $('<div class="card light card-played"><p class="text">' + text + '</p></div>');

          } else {
            card = $('<div class="card light"><p class="text">' + text + '</p><a href="#" data-payload="?card=' + id + '&player=' + localStorage.getItem("player-id") + '" class="play">Play</a></div>');
          }

          $("#cards").append(card);
        }
      }

    },
    this.GameAlreadyInProgress = function () {
      if (this.GameAlreadyInProgressBool == true) {
        return true;
      }
      $("#start-game").remove();
      $("#lobby").remove();
      $("#play-area").fadeIn();
      $("#cards").empty();
      this.RenderCards();

      this.GameAlreadyInProgressBool = true;
    },
    this.UpdateScores = function () {
      $this = this;
      $("#score tbody").empty()
      for (var i in $this.GameData.Players) {
        player = $this.GameData.Players[i];

        scores = $('<tr><td>' + player.Name.replace("[", "").replace("]", "") + '</td><td>' + player.Score + '</td></tr>');
        $("#score tbody").append(scores);
      }
    }
  }

  game = new Game();
  game.PollGame();

  function populateLobby() {
    $("#lobby #players").empty();
    for (var i in game.GameData.Players) {
      text = game.GameData.Players[i].Name;
      text = text.replace("[", "").replace("]", "");
      card = $('<div class="card light"><p class="text">' + text + ' has joined the game.</p></div>');
      $("#lobby #players").append(card);
    }
  }

  $("#lobby form#add-player").on('submit', function (e) {
    e.preventDefault();

    var name = $("#lobby form#add-player input#name").val();

    if (name.length == 0) {
      return false;
    }

    game.AddPlayer(name);

    window.setInterval(function () {
      populateLobby();
    }, 2000);

    $("#lobby form#add-player").fadeOut(1000);

  })

  $(".header #new-game").on('click', function (e) {
    e.preventDefault();
    game.New();
    location.reload();
  });
  $(".header #start-game").on('click', function (e) {
    e.preventDefault();
    game.Start();
    $("#lobby").remove();
    $("#play-area").fadeIn();
  })

  $("#cards").delegate('.card a.play', 'click', function (e) {
    e.preventDefault();
    d = $(this).data('payload')
    $this = this;

    $.getJSON( "/api/game/cards/play" + d, function( data ) {
      $($this).parent().addClass("card-played")
      $this.remove();
    });
  })

  $("#cards").delegate('.card a.award', 'click', function (e) {
    e.preventDefault();
    d = $(this).data('payload')
    $this = this;

    $.getJSON( "/api/game/cards/award" + d, function( data ) {
      e.preventDefault();
      $($this).parent().addClass("card-played")
      game.Start();
    });
  })
  
}( jQuery ));