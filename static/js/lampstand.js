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
	this.$slides = $(".slides");

}

PassageView.prototype.render = function(passages) {
  //Drop all sections with the lampstand class
  this.$slides.find(".lampstand").remove();
  //Recreate all the sections
  passages.forEach(function(el, index, array) {
  	this._renderPassage(el);
  }, this);
}

PassageView.prototype._renderPassage = function(passage) {
	var verses = passage.verses;
	var reference = passage.reference;
	var version = passage.version;
	var newPassage =  $("<section class='lampstand future'></section>").appendTo(this.$slides);
	var newSection;
	var fragments = 2;
	for(var i = 0; i < verses.length; i++) {
		var text = verses[i].text;
		var verseNo = verses[i].verseNo;
		if(i % fragments == 0) {
			newSection = $("<section class='lampstand'>").appendTo(newPassage);
		}
		newSection.append($("<span class='passage fragment'> <sup>" + verseNo +"</sup>" + text +"</span>"));
	}
	$(newPassage).find("section.lampstand").append($("<div class='reference'>" + reference + " (" + version + ")" + "</div>"));
}

/**
 * The SearchBar handles the actions of the searchbar, makes an AJAX call to the backend,
 * and passes the data to the PassageStore
 */
function SearchBar(passageStore) {
	this.passage = "";
	this.passageStore = passageStore;

	this.$version = $("#version");
	this.$passage = $("#passage");
	this.$searchBtn = $("#search");
	this.$clearAllBtn = $("#clear");
	this.$slides = $(".slides");

	this.$searchBtn.on("click", this.search.bind(this))
	this.$clearAllBtn.on("click", this.passageStore.deleteAll.bind(this.passageStore));
}

SearchBar.prototype.search = function() {
	console.log("Getting Data");
	var encodedVersion = encodeURI(this.$version.val());
	var encodedPassage = encodeURI(this.$passage.val());
	var pqUrl = "/api/" + encodedVersion + "/verses?passage=" + encodedPassage

	$.ajax({
		url: pqUrl,
		dataType: "json",
		success: this._dispatch.bind(this)
	})
}

//TODO: Handle errors
SearchBar.prototype._dispatch = function(data) {
	this.passageStore.addPassage(data);
}