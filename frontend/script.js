const EMPTY = 0;
const PLAYER_X = 1;
const PLAYER_O = 2;
let currentPlayer = PLAYER_X;
let gameOver = false;
let board;

function final() {
  return winner() !== EMPTY || !board.includes(EMPTY);
}

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
    input.addEventListener("click", renderBoard);
  });
}

function renderBoard() {
  if (currentPlayer === PLAYER_X) {
    this.value = "X";
    currentPlayer = PLAYER_O;
    aiMove()
  }else if(currentPlayer === PLAYER_O){
    this.value = "O";
    currentPlayer = PLAYER_X;
  }
  
  this.className = "inputOff";
  getBoardValues()
  board = nextTable(this)
  console.log(board)
}

function checkWinner(board) {
  const winPatterns = [
      [0, 1, 2], [3, 4, 5], [6, 7, 8], // Rows
      [0, 3, 6], [1, 4, 7], [2, 5, 8], // Columns
      [0, 4, 8], [2, 4, 6]             // Diagonals
  ];

  for (let pattern of winPatterns) {
      const [a, b, c] = pattern;
      if (board.cells[a] !== EMPTY && board.cells[a] === board.cells[b] && board.cells[a] === board.cells[c]) {
          return true;
      }
  }
  return false;
}

function playerMove(index) {
  if (!gameOver && board[index] === EMPTY) {
      board[index] = currentPlayer;
      renderBoard();
      if (final()) {
          gameOver = true;
          showResult();
      } else {
          currentPlayer = currentPlayer === PLAYER_X ? PLAYER_O : PLAYER_X;
          if (currentPlayer === PLAYER_O) {
              setTimeout(() => {
                  aiMove();
              }, 500); // Delay AI move for better UX
          }
      }
  }
}

function aiMove() {
  if (!gameOver) {
      let bestScore = -Infinity;
      let move;
      for (let i = 0; i < board.length; i++) {
          if (board[i] === EMPTY) {
              board[i] = PLAYER_O;
              let score = minimax(board, false);
              board[i] = EMPTY;
              if (score > bestScore) {
                  bestScore = score;
                  move = i;
              }
          }
      }
      board[move] = PLAYER_O;
      renderBoard();
      if (final()) {
          gameOver = true;
      } else {
          currentPlayer = PLAYER_X;
      }
  }
}

function minimax(board, isMaximizing) {
  let result = winner();
  if (result !== EMPTY) {
      return result === PLAYER_O ? 1 : -1;
  }
  if (!board.includes(EMPTY)) {
      return 0;
  }
  if (isMaximizing) {
      let bestScore = -Infinity;
      for (let i = 0; i < board.length; i++) {
          if (board[i] === EMPTY) {
              board[i] = PLAYER_O;
              let score = minimax(board, false);
              board[i] = EMPTY;
              bestScore = Math.max(score, bestScore);
          }
      }
      return bestScore;
  } else {
      let bestScore = Infinity;
      for (let i = 0; i < board.length; i++) {
          if (board[i] === EMPTY) {
              board[i] = PLAYER_X;
              let score = minimax(board, true);
              board[i] = EMPTY;
              bestScore = Math.min(score, bestScore);
          }
      }
      return bestScore;
  }
}


function nextTable(lastMove) {
  var focusTable = lastMove.id.slice(6) % 9;
  console.log(focusTable);
  var index = 0;

  document.querySelectorAll(".table").forEach((board) => {
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


//adaptar playermove para colocar no lugar de renderbord dentro 
//de inputs, transformar nextTable em um array global
//terminar de adaptar as funções