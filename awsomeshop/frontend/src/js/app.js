// Application entry point

// Test ES6+ features to verify Babel transpilation
const testArrowFunction = () => {
  console.log('Arrow function works');
};

const testAsyncAwait = async () => {
  const promise = Promise.resolve('Async/await works');
  const result = await promise;
  console.log(result);
};

const testSpreadOperator = () => {
  const arr = [1, 2, 3];
  const newArr = [...arr, 4, 5];
  console.log('Spread operator works:', newArr);
};

// Test optional chaining (ES2020)
const testOptionalChaining = () => {
  const obj = { a: { b: { c: 'Optional chaining works' } } };
  console.log(obj?.a?.b?.c);
};

// Initialize app
console.log('AWSomeShop Frontend - Babel transpilation test');
testArrowFunction();
testAsyncAwait();
testSpreadOperator();
testOptionalChaining();
