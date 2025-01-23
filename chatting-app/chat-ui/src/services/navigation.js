// navigation.js
let navigate;

export const setNavigate = (navigateFunction) => {
  navigate = navigateFunction;
};

export const getNavigate = () => {
  if (!navigate) {
    throw new Error("Navigate function has not been set. Call setNavigate first.");
  }
  return navigate;
};