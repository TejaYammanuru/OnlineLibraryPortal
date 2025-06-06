class ApplicationController < ActionController::API
  before_action :configure_permitted_parameters, if: :devise_controller?
  include Pundit::Authorization
  rescue_from Pundit::NotAuthorizedError, with: :user_not_authorized
  protected
  def configure_permitted_parameters
    devise_parameter_sanitizer.permit(:sign_up, keys: %i[name role])
    devise_parameter_sanitizer.permit(:account_update, keys: %i[name role])
  end
  private
  def user_not_authorized(exception)
    render json: { error: "Not authorized to perform this action" }, status: :forbidden
  end
end

