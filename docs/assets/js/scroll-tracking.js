// indexOf polyfill
var contains = function(needle) {
  // Per spec, the way to identify NaN is that it is not equal to itself
  var findNaN = needle !== needle;
  var indexOf;

  if(!findNaN && typeof Array.prototype.indexOf === 'function') {
    indexOf = Array.prototype.indexOf;
  } else {
    indexOf = function(needle) {
      var i = -1, index = -1;

      for(i = 0; i < this.length; i++) {
        var item = this[i];

        if((findNaN && item !== item) || item === needle) {
          index = i;
          break;
        }
      }

      return index;
    };
  }

  return indexOf.call(this, needle) > -1;
};

var sent = [];
var CONTACT = 'contact';
var marks = [10, 25, 50, 75, 100, CONTACT];
var contactForm = document.querySelector('#contact .container');

function isInViewport(elem) {
  var rect = elem.getBoundingClientRect();
  var html = document.documentElement;
  return (
    rect.top >= 0 &&
    rect.left >= 0 &&
    rect.bottom <= (window.innerHeight || html.clientHeight) &&
    rect.right <= (window.innerWidth || html.clientWidth)
  );
}

(function () {
  function sendScrollEvent() {
    var body = document.body;
    var html = document.documentElement;
    var docHeight = Math.max(body.scrollHeight, body.offsetHeight, html.clientHeight, html.scrollHeight, html.offsetHeight);
    var scrollHeight = document.body.scrollTop;
  
    var percentage = Math.floor((scrollHeight / docHeight) * 100);
    var alreadySent = contains.call(sent, percentage);
    var listeningFor = contains.call(marks, percentage);
    var contactFormAlreadySent = contains.call(sent, CONTACT);

    if (contactForm && !contactFormAlreadySent) {
      if (isInViewport(contactForm)) {
        sent.push(CONTACT);
        window.ga('send', 'event', 'Scroll', CONTACT, window.location.pathname);
      }
    }

    if (alreadySent || !listeningFor) return;
  
    sent.push(percentage);
  
    window.ga('send', 'event', 'Scroll', percentage + '%', window.location.pathname);
  }
  if (window.addEventListener) window.addEventListener('scroll', sendScrollEvent, false);
  else if (window.attachEvent) window.attachEvent('onscroll', sendScrollEvent);
})();
