function Del(elem) {
    var Reg = $(elem).data("id");
    var Button = document.getElementById("disable-"+Reg)
    console.log(Button.checked,Reg)
    var names = []
    var elm = document.getElementsByName(Reg);
    
    for (index = elm.length - 1; index >= 0; index--) {
        names.push("#"+elm[index].getAttribute('id'))
    }

    if(Button.checked){
      jQuery(names.toString()).fadeOut();

    } else {
      jQuery(names.toString()).fadeIn();
    }
  }

  $(function () {
    var load = function (url) {
        $.get(url).done(function (data) {
            $("html").html(data);
        })
    };

    $(document).on('click', 'a', function (e) {
        e.preventDefault();

        var $this = $(this),
            url = $this.attr("href"),
            title = $this.text();

            if (/(.*youtube.+)/.test(url)) {
              console.log("Youtube");
              window.open(
                url,
                '_blank'
              );
              return
            } else if (/(.*space\.bilibili.+)/.test(url)) {
              console.log("Bilibili");
              window.open(
                url,
                '_blank'
              );
              return
            } else if (url == "#") {
              return
            }

        history.pushState({
            url: url,
            title: title
        }, title, url);

        document.title = title;
        load(url, function(response){
          $("html").html(response);
        });
    });

    $(window).on('popstate', function (e) {
        var state = e.originalEvent.state;
        if (state !== null) {
            document.title = state.title;
            load(state.url, function(response){
              $("html").html(response);
            });
        } else {
            document.title = 'Go-Simp';
            if (state == null) {
              load("/", function(response){
                $("html").html(response);
              });
            } else {
              load(state.url, function(response){
                $("html").html(response);
              });
            }
        }
    });
});


! function(a) {
  a.fn.isOnScreen = function(b) {
      var c = this.outerHeight(),
          d = this.outerWidth();
      if (!d || !c) return !1;
      var e = a(window),
          f = {
              top: e.scrollTop(),
              left: e.scrollLeft()
          };
      f.right = f.left + e.width(), f.bottom = f.top + e.height();
      var g = this.offset();
      g.right = g.left + d, g.bottom = g.top + c;
      var h = {
          top: f.bottom - g.top,
          left: f.right - g.left,
          bottom: g.bottom - f.top,
          right: g.right - f.left
      };
      return "function" == typeof b ? b.call(this, h) : h.top > 0 && h.left > 0 && h.right > 0 && h.bottom > 0
  }
}(jQuery);

$(document).ready(function() {
  checkDisplay();

$(window).on('resize scroll', function() {
  checkDisplay();
});
});

function checkDisplay(){
$('.prg-count').each(function() {
    var $this = $(this);
    if ($this.isOnScreen()) {
      var countTo = $this.attr('data-count');
      $({
        countNum: $this.text()
      }).animate({
        countNum: countTo
      }, {
        duration: 8000,
        easing: 'linear',
        step: function() {
          $this.text(Math.floor(this.countNum));
        },
        complete: function() {
          $this.text(this.countNum);
          //alert('finished');
        }
      });
    }
  });
}

function IsEmpty() {
  if (document.forms['form'].Nickname.value === "" || document.forms['form'].Region.value === "") {
    alert("Nickname or Region empty");
    return false;
  }
  return true;
}
/*

  $('.GetMembers').on("click",function(){
    GroupName = $(this).attr("data_group");
    var urls = "/Group/"+GroupName
    $.ajax(
    {
        type:"GET",
        url: urls,
        success: function(response){
          document.title = response.pageTitle;
          window.history.pushState({"html":response.html,"pageTitle":response.pageTitle},"Go-Simp", urls);
          $("body").html(response);
        },
    })
  });

  $('.MemberInfo').on("click",function(){
    MemberName = $(this).attr("data_group");
    var urls = "/Member/"+MemberName
    $.ajax(
    {
        type:"GET",
        url: urls,
        success: function(response){
          document.title = response.pageTitle;
          window.history.pushState({"body":response.body,"pageTitle":response.pageTitle},MemberName, urls);
          $("body").html(response);
        },
    })
  });

  $('.navbar-brand').on("click",function(){
    $.ajax(
    {
        type:"GET",
        url: "/",
        success: function(response){
          document.title = response.pageTitle;
          window.history.pushState({"body":response.body,"pageTitle":response.pageTitle},"Go-Simp", "/");
          $("body").html(response);
        },
    })
  });

*/