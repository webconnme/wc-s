var initValue = function(name) {
  $.get("/module/" + name + "/properties", function(data) {
    for (var k in data) {
      $("#" + name + "-" + k).val(data[k]);
    }
  });
}

var saveProperties = function(module) {
  var params = $('#' + module + '-form').serializeArray();

  var data = {};
  data.name = module;
  data.properties = {};

  for (var i in params) {
    var name = params[i].name.replace(module + '-', "");
    if ($.isNumeric(params[i].value)) {
      data.properties[name] = parseInt(params[i].value);
    } else {
      data.properties[name] = params[i].value;
    }
  }

  console.log(JSON.stringify(data));

  $.post("/module/" + module + "/properties", JSON.stringify(data), function(data) {
    console.log(data);
  });
  return false;
}

