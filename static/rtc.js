let localVideo
let remoteVideo
let uuid
let serverConnection
let peerConnection
let peerConnectionConfig = {
  'iceServers': [
    {'url': 'stun:stun.stunprotocol.org:3478'},
    {'url': 'stun:stun.l.google.com:19302'}
  ],
  'turnServers': [
    {'url': 'turn:turn01.hubl.in?transport=udp'},
    {'url': 'turn:turn02.hubl.in?transport=tcp'}
  ]
}

const constraints = {
    video: true,
    audio: true
  }

uuid = createUUID();

function pageReady() {

  localVideo = document.getElementById('localVideo');
  remoteVideo = document.getElementById('remoteVideo');

  serverConnection = new WebSocket('wss://' + window.location.hostname + '/ws');
  serverConnection.onmessage = gotMessageFromServer;

  if(navigator.mediaDevices.getUserMedia) {
    navigator.mediaDevices.getUserMedia(constraints).then(getUserMediaSuccess).catch(err => {
      console.log(err);
    });
  } else {
    alert('NOT supported :(');
  }
}

function getUserMediaSuccess(stream) {
  localStream = stream;
  localVideo.srcObject = stream;
}

function gotMessageFromServer(message) {
  console.log('GOT MESSAGE FROM SERVER: ', message)
  if(!peerConnection) start(false);

   const signal = JSON.parse(message.data);
   console.log(signal)
   console.log("SDP?: ", signal.sdp)
   console.log("UUID: ", signal.uuid)

   // ignore self originated messages
   if(signal.uuid === uuid) return;

  if(signal.ice && signal.ice.candidate) {
    console.log("Adding ICE candidate:", signal.ice);
    peerConnection.addIceCandidate(new RTCIceCandidate(signal.ice)).catch(err => {
      console.log(err)
    })
  } else if(signal.sdp && signal.sdp.type) {
    peerConnection.setRemoteDescription(new RTCSessionDescription(signal.sdp)).then(() => {
      // only create answers in response to offers
      console.log('REMOTE DESCRIPTION SET');
      if(signal.sdp.type == 'offer') {
        peerConnection.createAnswer().then(createdDescription).catch(err => {
          console.log(err);
        })
      }
    }).catch(err => {
      console.log(err)
    })
  }
}

function gotIceCandidate(event) {
  if(event.candidate != null) {
    console.log('ice candidate received: ', event.candidate)
    serverConnection.send(JSON.stringify({'ice': event.candidate, 'uuid': uuid}))
  }
}

function createdDescription(description) {
  console.log('description received', description)

  peerConnection.setLocalDescription(description).then(() => {
    console.log('description set...')
    serverConnection.send(JSON.stringify({'sdp': peerConnection.localDescription, 'uuid': uuid}))
  }).catch(err => {
    console.log(err)
  })
}

function gotRemoteStream(event) {
  console.log('got remote stream', event.streams);
  remoteVideo.srcObject = event.streams[0];
}

function start(isCaller) {
  peerConnection = new RTCPeerConnection(peerConnectionConfig);
  peerConnection.onicecandidate = gotIceCandidate;
  peerConnection.ontrack = gotRemoteStream;
  peerConnection.addStream(localStream);

  if(isCaller) {
    peerConnection.createOffer().then(
      createdDescription
    ).catch(err => {
      console.log(err)
    })
  }
}

// https://stackoverflow.com/a/105074/515584
function createUUID() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x1000).toString(16).substring(1);
  }

  return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4()
}
