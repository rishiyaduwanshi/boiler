// __author Abhinav Prakash
// __desc Async error handler
// __version 1
// __var bl__CLASS_NAME = AppError
// __var bl__FUNCTION_NAME, bl__ERROR_DESC, 
// __var bl__ERROR_MESSAGE = "Internal Server Error"

class bl__CLASS_NAME extends Error {
    constructor({ statusCode = 500, message = "bl_ERROR_MESSAGE", errors = [], stack = "" }) {
      super(message);
      this.statusCode = statusCode;
      this.success = false;
      this.errors = errors;
  
      if (stack) {
        this.stack = stack;
      } else {
        Error.captureStackTrace(this, this.constructor);
      }
    }
  }
  
  export { bl__CLASS_NAME };
  


export class NotFoundError extends AppError {
  constructor(message = 'Resource not found') {
    super({ message, statusCode: 404 });
  }
}

export class BadRequestError extends AppError {
  constructor(message = 'Bad request') {
    super({ message, statusCode: 400 });
  }
}

export class UnauthorizedError extends AppError {
  constructor(message = 'Unauthorized') {
    super({ message, statusCode: 401 });
  }
}
