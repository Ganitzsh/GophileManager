var busy = 0
var id = 0
var socket

window.onbeforeunload = function (e) {
	if (busy != 0) {
		return 'Operations in progress'
	}
	return nil
};


$(document)
	.ready(function () {
		$('.sk-fading-circle')
			.hide()
		fetchContent('current')
	});

function generateProgressAlert(text) {
	return `
	<div class="alert alert-info alert-dismissible" role="alert" id="alert-${id}">
		<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
		${text}
		<div class="progress">
				<div class="progress-bar progress-bar-striped active" role="progressbar" style="width: 100%">
				</div>
		</div>
	</div>
	`
}

function generateAlert(type, text) {
	return `
	<div class="alert alert-${type} alert-dismissible" role="alert" id="alert-${id}">
		<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
		${text}
	</div>
	`
}

function newNotification(type, message) {
	var alert
	switch (type) {
	case 'progress':
		alert = generateProgressAlert(message)
		break;
	case 'done':
		alert = generateDoneAlert(message)
		break;
	}
	$('#notifications')
		.append(alert)
	id++
	return '#alert-' + (id - 1)
		.toString()
}

function loading() {
	$('#files')
		.hide()
	$('.sk-fading-circle')
		.show();
}

function doneLoading() {
	$('#files')
		.show()
	$('.sk-fading-circle')
		.hide();
}

function fetchContent(where) {
	busy++
	var str = '/app/files/' + where
	var uri = encodeURI(str)
	loading()
	fetch(uri, {
			credentials: 'same-origin' // Add cookie
		})
		.then(function (response) {
			if (response.ok) {
				response.text()
					.then(function (html) {
						$('html, body')
							.animate({
								scrollTop: 0
							}, 'fast');
						$('#files')
							.html(html)
						$('[data-toggle="tooltip"]')
							.tooltip()
						doneLoading()
						busy--
					});
			} else {}
		})
		.catch(function (error) {});
}

function watch(name) {
	window.location = "/app/watch/" + name;
}

function convertToMP4(name) {
	busy++
	$('#progressModal')
		.modal('show')
	if (busy) {
		console.log("Busy");
		$('#progressModal')
			.modal('hide')
		return
	}
	var str = '/app/convert/' + name
	var uri = encodeURI(str)
	fetch(uri, {
			credentials: 'same-origin'
		})
		.then(function (response) {
			if (response.ok) {
				response.json()
					.then(function (json) {});
				fetchContent('current')
				busy--
			} else {
				$('#progressModal')
					.modal('hide')
			}
		})
		.catch(function (error) {
			$('#progressModal')
				.modal('hide')
		});
}

function deleteFile(name) {
	bootbox.confirm("Delete <strong>" + name + "</strong>?", function (result) {
		if (result) {
			busy++
			var str = '/app/delete/' + name
			var uri = encodeURI(str)
			fetch(uri, {
					credentials: 'same-origin'
				})
				.then(function (response) {
					if (response.ok) {
						fetchContent('current')
					} else {}
				})
				.catch(function (error) {});
		}
	});
}

function setToCompress(name) {
	$('#fileName')
		.text(name)
	$('#myModal')
		.modal('show')
}

function validateCompressionForm() {
	var name = $('#archiveFileName')
		.val()
	console.log(name);
	if (!name) {
		$('#compressionFormGroup')
			.addClass('has-error')
		$('#nameHelp')
			.removeClass('hidden')
	} else {
		var file = $('#fileName')
			.text()
		compressAndDownload(name, file)
	}
}

function downloadFile(name) {
	busy++
	var str = '/app/download/' + name
	var uri = encodeURI(str)
	console.log(uri);
	fetch(uri, {
			credentials: 'same-origin'
		})
		.then(function (response) {
			if (response.ok) {
				var mime = response.headers.get('Content-Type')
				response.blob()
					.then(function (blob) {
						download(blob, name, mime)
						busy--
					});
			} else {}
		})
		.catch(function (error) {
			console.log('There has been a problem with your fetch operation: ' + error.message);
		});
}

function compressAndDownload(name, target) {
	var alertID = newNotification('progress', 'Compressing <strong>' + target + '</strong>')
	console.log(alertID);
	$('#fileName')
		.text('')
	$('#archiveFileName')
		.val('')
	$('#myModal')
		.modal('hide')
	console.log(name + ' - ' + target);
	var str = '/app/compress/' + target + '/name/' + name
	var uri = encodeURI(str)
	console.log(uri);
	fetch(uri, {
			credentials: 'same-origin'
		})
		.then(function (response) {
			if (response.ok) {
				fetchContent('current')
				$(alertID)
					.replaceWith(generateAlert('success', `<strong>${target}</strong> commpressed successfully!`))
			} else {
				response.json()
					.then(function (error) {
						console.log(error);
						$(alertID)
							.replaceWith(generateAlert('danger',
								`<p>Compression failed for <strong>${target}</strong></p>
								<p><strong>Cause:</strong> ${error.message}</p>`
							))
					})
			}
		})
		.catch(function (error) {
			$('#progressModal')
				.modal('hide')
			console.log('There has been a problem with your fetch operation: ' + error.message);
		});
}
