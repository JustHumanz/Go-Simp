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


const counters = document.querySelectorAll('.counter');
const speed = 200; // The lower the slower

counters.forEach(counter => {
const updateCount = () => {
const target = +counter.getAttribute('data-target');
const count = +counter.innerText;

// Lower inc to slow and higher to slow
const inc = target / speed;

// console.log(inc);
// console.log(count);

// Check if target is reached
if (count < target) {
// Add inc to count and output in counter
counter.innerText = count + inc;
// Call function every ms
setTimeout(updateCount, 1);
} else {
counter.innerText = target;
}
};

updateCount();
});

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