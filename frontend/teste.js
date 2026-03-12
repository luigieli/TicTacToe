var player = 1;
var bigTable;

function createBoard() {
  var contInput = 0;
  var board_container = document.getElementById("board-container");
  bigTable = ["", "", "", 
              "", "", "", 
              "", "", ""];
              
  for (var i = 0; i < 9; i++) {
    var divTable = document.createElement("span");
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

function createLog(message){
  var logDiv = document.getElementById('log');
  var existingParagraphs = logDiv.getElementsByTagName('p');

  for (var i = 0; i < existingParagraphs.length; i++) {
    existingParagraphs[i].style.animation = 'slideUpOut 1s ease-in forwards';
    existingParagraphs[i].style.animation = 'vanish 1s ease-in forwards';
  }

  var paragrafo = document.createElement('p');
  paragrafo.textContent = message;
  
  logDiv.appendChild(paragrafo);
  logDiv.style.overflow = "visible"
  
}

// setInterval(() => {
//   createLog('Nova mensagem: ' + new Date().toLocaleTimeString());
// }, 5000);

function interaction() {
  checkPlayer(this);
  checkWinner(this);
  nextTable(this);
  getBoardValues();
  verifyBigTable();
}

function gameOver(){
  const allTables = [...document.querySelectorAll('span.table')];
  allTables.forEach((table) => {
    table.className = "tableOff"
  })
}

function verifyBigTable() {
  var bigTableChanged = switchValueZ()
  if(bigTableChanged != null && bigTableChanged != "TIEX") {
    console.log("O grande vencedor do tabuleiro é " + bigTableChanged)
    createLog('O grande vencedor do tabuleiro é ' + bigTableChanged)
    gameOver()
  }
  else if(checkTie(bigTable) || bigTableChanged == "TIEX"){
    console.log("A grande tabela nao possui vencedor")
    createLog('A grande tabela nao possui vencedor')
    gameOver()
  } 
}

function switchValueZ() {
  var alterBigTableX = bigTable.slice()
  var alterBigTableO = bigTable.slice()
  
  for(let i = 0; i < 8; i++){
    if(alterBigTableO[i] === 'Z')
      alterBigTableO[i] = 'O'
  }

  for(let i = 0; i < 8; i++){
    if(alterBigTableX[i] === 'Z')
      alterBigTableX[i] = 'X'
  }
  
  if(checkWinnerInMatrix(alterBigTableO) == checkWinnerInMatrix(alterBigTableX)) return checkWinnerInMatrix(bigTable)
  else if(checkWinnerInMatrix(alterBigTableO) && checkWinnerInMatrix(alterBigTableX)) return "TIEX"
  else if(checkWinnerInMatrix(alterBigTableO)) return 'O'
  else if(checkWinnerInMatrix(alterBigTableX)) return 'X'
  return null
}

function checkWinner(inputClicked) {
  var tableHTMLClicked = inputClicked.parentElement.parentElement;
  var matrix = [];

  let inputs = tableHTMLClicked.querySelectorAll('input[type="button"]');
  inputs.forEach((input) => {
    matrix.push(input.value);
  });

  var winner = checkWinnerInMatrix(matrix);
  var boards = document.querySelectorAll("span");
  var indiceTable = Array.from(boards).indexOf(tableHTMLClicked);

  if (winner !== null) {
    console.log("O jogador " + winner + " venceu a tabela " + (indiceTable+1) + "!")
    createLog("O jogador " + winner + " venceu a tabela " + (indiceTable+1) + "!")
    bigTable[indiceTable] = winner;
    tableHTMLClicked.className = "tableOff";
  } else {
    var tie = checkTie(matrix);
    if (tie) {
      console.log("A tabela " + indiceTable + " empatou! ");
      createLog("A tabela " + indiceTable + " empatou! ")
      tableHTMLClicked.className = "tableOff";
      bigTable[indiceTable] = 'Z';
    }
  }
}

function checkWinnerInMatrix(matrix) {
  for (let i = 0; i < 3; i++) {
    if (
      matrix[i * 3] === matrix[i * 3 + 1] &&
      matrix[i * 3] === matrix[i * 3 + 2] &&
      matrix[i * 3] !== ""
    ) {
      return matrix[i * 3];
    }
  }

  for (let i = 0; i < 3; i++) {
    if (
      matrix[i] === matrix[i + 3] &&
      matrix[i] === matrix[i + 6] &&
      matrix[i] !== ""
    ) {
      return matrix[i];
    }
  }

  if (
    (matrix[0] === matrix[4] && matrix[0] === matrix[8] && matrix[0] !== "") ||
    (matrix[2] === matrix[4] && matrix[2] === matrix[6] && matrix[2] !== "")
  ) {
    return matrix[4];
  }

  return null;
}

function checkTie(matrix) {
  var confirmation = true;
  matrix.forEach((cell) => {
    if (cell === "") {
      confirmation = false;
    }
  });
  return confirmation;
}

function checkPlayer(thisInput) {
  thisInput.value = player % 2 != 0 ? "X" : "O";
  thisInput.className = "inputOff";
  player++;
}

function nextTable(lastMove) {
  var focusTable = lastMove.id.slice(6) % 9;
  var index = 0;
  var boards = document.querySelectorAll("span");

  if (boards[focusTable].className == "tableOff") {
    console.log("Este tabuleiro nao pode ser escolhido!");
    boards.forEach((board) => {
      board.style.pointerEvents = "auto";
      board.style.opacity = 1;
      board.style.cursor = "auto";
    });
    return getBoardValues()[focusTable];
  }

  boards.forEach((board) => {
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
  var boards = document.querySelectorAll("span");

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
  createBoard()
} 

startGame();
