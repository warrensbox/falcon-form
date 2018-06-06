(function() {
  'use strict';
  window.addEventListener('load', function() {
    var form = document.getElementById('eloquaForm');
    if (!form) return;
    form.addEventListener('submit', function(event) {
      // Prevent form submission and contact with server
      if (form.checkValidity() == false) {
        event.preventDefault();
        event.stopPropagation();
      }

      // Add a class when we attempt to validate so we can scope CSS to this state.
      // Otherwise, browsers load with the `:invalid` state applied.
      form.classList.add('validated');
    }, false);
  }, false);
}());


//
// Mobile nav
//
var toggle = document.querySelector('.site-header-toggle');
var siteHeader = document.querySelector('.site-header')
toggle.addEventListener('click', function() {
  siteHeader.classList.toggle('open');
  toggle.setAttribute('aria-expanded', !toggle.getAttribute('aria-expanded'));
});


//
// Show more (than just 9) event cards
//
var eventCards = document.querySelector('.event-cards');
var moreEventsBtn = document.querySelector('.js-more-events');
if (eventCards && moreEventsBtn) {
  var events = eventCards.querySelectorAll('.js-event-card');
  moreEventsBtn.addEventListener('click', function() {
    var hiddenCards = eventCards.querySelectorAll('.js-hidden-card');

    // Determine number to start showing cards
    var numVisible = events.length - hiddenCards.length;

    // Starting at numVisible, show the next 9 cards
    for (var i = numVisible; i < numVisible + 9; i++) {
      if (events[i]) {
        events[i].classList.remove('js-hidden-card');
      } else {
        // Hide the button if there are no more events
        moreEventsBtn.classList.add('d-none');

        // End the loop
        break;
      }
    }
  });
}


//
// Track click events in Google Analytics
//
var gaEvents = document.querySelectorAll('[data-ga-click]');
for (var i = 0; i < gaEvents.length; i++) {
  var el = gaEvents[i];
  el.addEventListener('click', function(e) {
    var gaClick = e.target.dataset.gaClick;
    var events = gaClick.split(', ');
    var category = events[0];
    var action = events[1];
    var label = events[2];
    
    window.ga('send', 'event', category, action, label);
  });
}


//
// GeoIP Script
//
(function() {
  "use strict";
  window.addEventListener("load", function() {
    var onSuccess = function(location) {
      var country = location.country.iso_code;
      var countryEl = document.getElementById("country");
      if (countryEl) countryEl.value = country;

      var subdivision = location.most_specific_subdivision.names.en;
      var subdivisionEl = document.getElementById("subdivision")
      if (subdivisionEl) subdivisionEl.value = subdivision;

      var city = location.city.names.en;
      var cityEl = document.getElementById("city")
      if (cityEl) cityEl.value = city;
    }

    geoip2.city(onSuccess);
  }, false);
}());


//
// Scroll Anchors
//
var anchors = document.querySelectorAll('.js-anchor');
for (var i = 0; i < anchors.length; i++) {
  var el = anchors[i];
  el.addEventListener('click', function(e) {
    var href = e.target.getAttribute('href').split('#')[1];
    var bits = href.split('#');
    var hash = bits[bits.length - 1];

    var section = document.getElementById(hash);
    if (section) {
      e.preventDefault();
      var yPos = section.getBoundingClientRect().top;
      scrollIt(
        section,
        300
      );

      // Adds the hash to the window's url; can
      // be removed if deemed unnecessary
      history.replaceState(null, null, '#' + hash);
    }
  });
}
