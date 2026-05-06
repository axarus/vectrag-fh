interface ErrorAlertProps {
    message: string;
    onRetry?: () => void;
  }
  
  export default function ErrorAlert({ message, onRetry }: ErrorAlertProps) {
    return (
      <div className="alert alert-error">
        <span>{message}</span>
        {onRetry && (
          <button className="btn btn-sm" onClick={onRetry}>
            Retry
          </button>
        )}
      </div>
    );
  }
  