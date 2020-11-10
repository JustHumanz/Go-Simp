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
  var counter = function(data) {
    data.ready(function() {
      if ($('h2.counter').length) {
        $('.counter').each(function() {
          var $this = $(this),
              countTo = $this.attr('data-count');
          
          $({ countNum: $this.text()}).animate({
            countNum: countTo 
          },
          {
            duration: 8000,
            easing:'linear',
            step: function() {
              $this.text(Math.floor(this.countNum));
            },
            complete: function() {
              $this.text(this.countNum);
            }
          });  
        });
      }
    });
  }
  
  counter($('body'))

  $(document).on('click', 'a', function (e) {
      e.preventDefault();

      var $this = $(this),
          url = $this.attr("href"),
          title = $this.text();

          if (url == "#") {
            return
          } else if (/(.*youtube.+)/.test(url)) {
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
          }

      history.pushState({
          url: url,
          title: title
      }, title, url);

      document.title = title;
      $('body').load(url,function(){
        counter($('body'))
      });
  });

  $(window).on('popstate', function (e) {
      var state = e.originalEvent.state;
      if (state !== null) {
          document.title = state.title;
          $('body').load(state.url,function(){
            counter($('body'))
          });
      } else {
          document.title = 'Go-Simp';
          if (state == null) {
            $('body').load("/",function(){
              counter($('body'))
            });
          } else {
            $('body').load(state.url,function(){
              counter($('body'))
            });
          }
      }
  });
});

function IsEmpty() {
if (document.forms['form'].Nickname.value === "" || document.forms['form'].Region.value === "") {
  alert("Nickname or Region empty");
  return false;
}
return true;
}



/*
counter = document.getElementsByClassName("counter")
for (index = counter.length - 1; index >= 0; index--) {
    console.log(counter[index].getAttribute('data-count'))


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