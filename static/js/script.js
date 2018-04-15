  $("#searchForm").on("submit", function (e) {
      e.preventDefault();
      $("#search").prop('disabled', true);
      var urlData = $("#url").val()
      if ($('#resultdiv').length){
          $('#resultdiv').remove();
      }
      console.log(urlData);
      $.post('/', {
          'url' : urlData
      }).done(function(data){
          var html="";
          $.each(data, function(index, item) {
          html=html+"<div style=\"display:inline-block; border: 1px solid #ccc; padding: 3px 10px; margin: 5px;width:20%;\"> "+index+" : <b>"+ item+"</b></div>";
          });
          $("#result").prepend("<div id=\"resultdiv\">The search result for "+urlData+"<br> "+html+"</div>")
          $("#search").prop('disabled', false);

      }).fail(function(xhr, status, error) {
            alert(xhr.responseText);
            $("#search").prop('disabled', false);
      });
      
 });
    $("#finderForm").on("submit", function (e) {
        e.preventDefault();
        $("#submit").prop('disabled', true);
        var start = $("#start").val()
        var row = $("#row").val()
        var column = $("#column").val()
        if ($('#resultdiv').length){
            $('#resultdiv').remove();
        }
        $.post('/finder', {
            'start' : start,
            'row' : row,
            'column' : column
        }).done(function(res){ 
            var data = $.parseJSON(res);
            console.log(data)
            var html="<table style=\"width:100%; border: 1px solid black;\">";
            var k=0;
            for (var i = 1; i <= column; i++) {
                html=html+"<tr>";
                for (var j = 1; j <= row; j++) {
                html=html+"<td style=\"border: 1px solid black;text-align: center;\">"+data.data[k++]+"</td>";
                }
                html=html+"</tr>";            
            }
            html=html+"</table>"
            $("#result").prepend("<div id=\"resultdiv\">The result is :"+html+"</div>");
            $("#submit").prop('disabled', false);
         }).fail(function(xhr, status, error) {
            $("#submit").prop('disabled', false);
            alert(xhr.responseText);
        });

    });


