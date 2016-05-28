// select a list of matching elements, context is optional
function qs(selector, context) {
    return (context || document).querySelectorAll(selector);
}

function qs1(selector, context) {
    return (context || document).querySelector(selector);
}

function getCORS(url, success) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', url);
    xhr.onload = success;
    xhr.send();
    return xhr;
}


/**
 * The passage store saves the passage items to local storage
 *
 */
function PassageStore(passageView) {
	this.loadStore();
	window.addEventListener("storage", this.onStorageChange.bind(this), false);
}

PassageStore.prototype.addPassage = function(passage) {
	this.passages.push(passage);
	this.saveStore();
}

PassageStore.prototype.deleteAll = function() {
	this.passages.length = 0;
	this.saveStore();
	Reveal.slide(0);
}

PassageStore.prototype.deletePassage = function(index) {
	this.passages.remove(index);
	this.saveStore();
}

PassageStore.prototype.saveStore = function() {
	localStorage["lampstand"] = JSON.stringify(this.passages);
	this.passageView.render(this.passages);
	Reveal.sync();
}


PassageStore.prototype.loadStore = function() {
	var passages = localStorage["lampstand"];
	if(typeof passages != "string") {
		localStorage["lampstand"] = "[]";
	}
	//try catch
	try {
		this.passages = JSON.parse(localStorage["lampstand"]);
	} catch(err) {
		this.passages = [];
	}
	this.passageView = passageView;
	this.passageView.render(this.passages);	
}

PassageStore.prototype.onStorageChange = function(e) {
	console.log("Storage Event Fired!");
	if(e.key == "lampstand") {
		this.passages = JSON.parse(e.newValue);
		this.passageView.render(this.passages);
		Reveal.sync();
	}
}

/**
 * The PassageView takes changes from the PassageStore and renders them
 * in the RevealJS section
 */
function PassageView() {
	this.$slides = qs1(".slides");

}

PassageView.prototype.render = function(passages) {
	//Drop all sections with the lampstand class
	var sections = qs(".lampstand")
	for(var i = 0; i < sections.length; i++) {
		sections[i].parentNode.removeChild(sections[i])
	}
	//Recreate all the sections
	passages.forEach(function(el, index, array) {
		this._renderPassage(el);
	}, this);
}

PassageView.prototype._renderPassage = function(passage) {
	var verses = passage.verses;
	var reference = passage.reference;
	var version = passage.version;
	var newPassage = document.createElement("section");
	newPassage.className = 'lampstand future';
	this.$slides.appendChild(newPassage);
	var newSection;
	var fragments = 2;
	for(var i = 0; i < verses.length; i++) {
		var text = verses[i].text;
		var verseNo = verses[i].verseNo;
		if(i % fragments == 0) {
			newSection = this._createNewSection()
			newVerse = this._createNewVerseEl(verseNo, text, false)
			newPassage.appendChild(newSection)
			newSection.appendChild(newVerse)
		} else {
			newVerse =  this._createNewVerseEl(verseNo, text, true)
			newSection.appendChild(newVerse)
		}
	}
	var passages = qs("section.lampstand", newPassage);
	var referenceEl = document.createElement("div");
	referenceEl.className = 'reference';
	referenceEl.innerHTML = reference + " (" + version + ")";
	for(var i = 0; i < passages.length; i++) {
		passages[i].appendChild(referenceEl.cloneNode(true));
	}
}

PassageView.prototype._createNewSection = function() {
	var newSection = document.createElement("section");
	newSection.className = 'lampstand part';
	return newSection;
}

PassageView.prototype._createNewVerseEl = function(verseNo, text, isFragment) {
	var newVerse = document.createElement("span");
	newVerse.className = newVerse.className = "passage"
	if(isFragment) {
		newVerse.className = newVerse.className + " fragment";
	}
	newVerse.innerHTML = "<sup>" + verseNo +"</sup>" + text;
	return newVerse;
}

/**
 * The SearchBar handles the actions of the searchbar, makes an AJAX call to the backend,
 * and passes the data to the PassageStore
 */
function SearchBar(passageStore) {
	this.passage = "";
	this.passageStore = passageStore;

	this.$version = qs1("#version");
	this.$passage = qs1("#passage");
	this.$searchBtn = qs1("#search");
	this.$clearAllBtn = qs1("#clear");
	this.$slides = qs1(".slides");

	this.$searchBtn.addEventListener("click", this.search.bind(this))	
	this.$clearAllBtn.addEventListener("click", this.passageStore.deleteAll.bind(this.passageStore));
}

SearchBar.prototype.search = function() {
	console.log("Getting Data");
	var encodedVersion = encodeURI(this.$version.value);
	var encodedPassage = encodeURI(this.$passage.value);
	var pqUrl = "/api/" + encodedVersion + "/verses?passage=" + encodedPassage;


	getCORS(pqUrl, this._dispatch.bind(this));

}

//TODO: Handle errors
SearchBar.prototype._dispatch = function(request) {
    var response = JSON.parse(request.currentTarget.response || request.target.responseText);
	this.passageStore.addPassage(response);
}