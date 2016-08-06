$(document).ready(function() {
  $('[data-toggle="tooltip"]').tooltip()
});

function watch(name) {
  window.location = "/app/watch/" + name;
}

function convertToMP4(name) {
  var str = '/app/convert/' + name
  var uri = encodeURI(str)
  fetch(uri).then(function(response) {
          if (response.ok) {
              response.json().then(function(json) {
                  console.log(json);
              });
              location.reload();
          } else {
              console.log('Network response was not ok.');
          }
      })
      .catch(function(error) {
          console.log('There has been a problem with your fetch operation: ' + error.message);
      });
}

function deleteFile(name) {
    var str = '/app/delete/' + name
    var uri = encodeURI(str)
    fetch(uri).then(function(response) {
            if (response.ok) {
                response.json().then(function(json) {
                    console.log(json);
                });
                location.reload();
            } else {
                console.log('Network response was not ok.');
            }
        })
        .catch(function(error) {
            console.log('There has been a problem with your fetch operation: ' + error.message);
        });
}

function setToCompress(name) {
    $('#fileName').text(name)
    $('#myModal').modal('show')
}

function validateCompressionForm() {
    var name = $('#archiveFileName').val()
    console.log(name);
    if (!name) {
        $('#compressionFormGroup').addClass('has-error')
        $('#nameHelp').removeClass('hidden')
    } else {
        var file = $('#fileName').text()
        compressAndDownload(name, file)
    }
}

function downloadFile(name) {
    var str = '/app/download/' + name
    var uri = encodeURI(str)
    console.log(uri);
    fetch(uri).then(function(response) {
            if (response.ok) {
                var mime = response.headers.get('Content-Type')
                response.blob().then(function(blob) {
                    download(blob, name, mime)
                    console.log(name);
                });
            } else {
                response.json().then(function(json) {
                    console.log(json);
                });
                console.log('Network response was not ok.');
            }
        })
        .catch(function(error) {
            console.log('There has been a problem with your fetch operation: ' + error.message);
        });
}

function compressAndDownload(name, target) {
    $('#fileName').text('')
    $('#archiveFileName').val('')
    $('#myModal').modal('hide')
    console.log(name + ' - ' + target);
    var str = '/app/compress/' + target + '/name/' + name
    var uri = encodeURI(str)
    console.log(uri);
    fetch(uri).then(function(response) {
            if (response.ok) {
                response.blob().then(function(myBlob) {
                    // download(myBlob, name + '.tar', 'application/x-tar');
                    location.reload();
                });
            } else {
                console.log('Network response was not ok.');
            }
        })
        .catch(function(error) {
            console.log('There has been a problem with your fetch operation: ' + error.message);
        });
}
