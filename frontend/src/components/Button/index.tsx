import Musa from "mini-framework";

function Button({ ...props }) {
  var { className, children, ...props } = props;
  
  return (
    <button 
    className={className ? "Button " + className : "Button"} 
    {...props}>
      {children}
    </button>
  );
}

export default Button;
