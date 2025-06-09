class Users::SessionsController < Devise::SessionsController
  include ActionController::MimeResponds
  respond_to :json

  # POST /login
  def create
    begin
      self.resource = warden.authenticate!(auth_options)
      sign_in(resource_name, resource)

      render json: {
        message: 'Logged in successfully',
        user: user_response(resource)
      }, status: :ok
    rescue => e
      render json: {
        message: 'Invalid email or password'
      }, status: :unauthorized
    end
  end

  # DELETE /logout
  def destroy
    sign_out(resource_name)
    render json: { message: 'Logged out successfully' }, status: :ok
  end

  # POST /verify_token
  def verify_token
    token = request.headers['Authorization']&.split(' ')&.last

    if token.blank?
      render json: { error: 'Token missing' }, status: :unauthorized and return
    end

    begin
      secret_key = 'Teja'
      decoded_token = JWT.decode(token, secret_key, true, algorithm: 'HS256')
      payload = decoded_token.first

      user_id = payload['sub'] || payload['id']
      user = User.find_by(id: user_id)

      if user
        render json: {
          id: user.id,
          email: user.email,
          name: user.name,
          role: user[:role]  
        }, status: :ok
      else
        render json: { error: 'User not found' }, status: :unauthorized
    end

    rescue JWT::DecodeError => e
      render json: { error: "Invalid token: #{e.message}" }, status: :unauthorized
    end
  end




  private

  def user_response(user)
    {
      id: user.id,
      email: user.email,
      name: user.name
    }
  end
end
