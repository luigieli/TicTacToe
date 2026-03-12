var player = 1;

function createBoard() {
  var contInput = 0;
  var board_container = document.getElementById("board-container");
  for (var i = 0; i < 9; i++) {
    var divTable = document.createElement("div");
    divTable.className = "table";

    for (var j = 0; j < 3; j++) {
      var divLine = document.createElement("div");
      for (var k = 0; k < 3; k++) {
        var divCell = document.createElement("div");
        divCell.className = "cell";
        var input = document.createElement("input");
        input.type = "button";
        input.id = `input-${contInput}`;
        contInput++;
        divCell.appendChild(input);
        divTable.appendChild(divCell);
      }
    }
    board_container.appendChild(divTable);
  }

  var inputs = document.querySelectorAll("input");
  var i = 0;
  inputs.forEach((input) => {
    input.value = "";
  });

  inputs.forEach((input) => {
    input.addEventListener("click", interaction);
  });
}

function interaction() {
  if (this.value === "") {
    this.value = player % 2 != 0 ? "X" : "O";
    this.className = "inputOff";
    player++;
  }
  
  getBoardValues()
  console.log(nextTable(this))
}

function nextTable(lastMove) {
  var focusTable = lastMove.id.slice(6) % 9;
  console.log(focusTable);
  var index = 0;

  document.querySelectorAll(".table").forEach((board) => {
    console.log(index);
    if (index != focusTable) {
      board.style.pointerEvents = "none";
      board.style.opacity = 0.5;
      board.style.cursor = "not-allowed";
    }

    if (index == focusTable && board.style.cursor == "not-allowed") {
      board.style.pointerEvents = "auto";
      board.style.opacity = 1;
      board.style.cursor = "auto";
    }

    index++;
  });

  return getBoardValues()[focusTable];
}

function getBoardValues() {
  var boardValues = [];
  var boards = document.querySelectorAll(".table");

  boards.forEach((board) => {
    var boardMatrix = [];
    var cells = board.querySelectorAll(".cell input");
    cells.forEach((cell, index) => {
      if (index % 3 === 0) boardMatrix.push([]);
      boardMatrix[boardMatrix.length - 1].push(cell.value);
    });
    boardValues.push(boardMatrix);
  });

  return boardValues;
}

function startGame() {
  createBoard();
}

startGame();
