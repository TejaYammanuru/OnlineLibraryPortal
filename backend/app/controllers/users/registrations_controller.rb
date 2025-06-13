class Users::RegistrationsController < Devise::RegistrationsController
  include ActionController::Flash  

  respond_to :json
  skip_before_action :authenticate_scope!, only: [:update]
  before_action :authenticate_user!, only: [:update, :destroy]

  def create
    Rails.logger.info "Create action started with params: #{sign_up_params.inspect}"
    build_resource(sign_up_params)
    

    if resource.save
      Rails.logger.info "User saved successfully: #{resource.inspect}"

      if current_user&.admin? && resource.librarian?
        LibrarianLog.create!(
          action: "create",
          admin: {
            id: current_user.id,
            name: current_user.name,
            email: current_user.email
          },
          librarian: {
            id: resource.id,
            name: resource.name,
            email: resource.email
          }
        )
      elsif resource.member?
        MemberLog.create!(
          action: "signup",
          member: {
            id: resource.id,
            name: resource.name,
            email: resource.email
          }
        )
        UserMailer.welcome_email(resource).deliver_later
      end

      if resource.active_for_authentication?
        sign_up(resource_name, resource)
        render json: { message: 'User registered successfully', user: user_response(resource) }, status: :created
      else
        expire_data_after_sign_in!
        render json: { message: 'Signed up but account is inactive', user: user_response(resource) }, status: :ok
      end
    else
      Rails.logger.info "User registration failed: #{resource.errors.full_messages}"
      render json: { message: 'User registration failed', errors: resource.errors.full_messages }, status: :unprocessable_entity
    end
  end

  def update
    Rails.logger.info "Update action started by user #{current_user&.id} with params: #{params.inspect}"

    if current_user.admin? && params[:id].present?
    
      Rails.logger.info "Admin updating librarian with id: #{params[:id]}"
      target_user = User.find_by(id: params[:id], role: User.roles[:librarian])

      unless target_user
       Rails.logger.info "Librarian not found with id: #{params[:id]}"
        return render json: { message: 'Librarian not found' }, status: :not_found
      end

      authorize target_user

      before_attrs = target_user.slice(:name, :email)

      if target_user.update(account_update_params)
        Rails.logger.info "Librarian updated successfully: #{target_user.inspect}"

        after_attrs = target_user.slice(:name, :email)

        LibrarianLog.create!(
          action: "update",
          admin: {
           id: current_user.id,
            name: current_user.name,
            email: current_user.email
          },
          librarian: {
            id: target_user.id,
           name: target_user.name,
            email: target_user.email
          },
          before_update: before_attrs,
          after_update: after_attrs
        )

        render json: { message: 'Librarian updated successfully', user: user_response(target_user) }, status: :ok
      else
        Rails.logger.info "Update failed for librarian: #{target_user.errors.full_messages}"
        render json: { message: 'Update failed', errors: target_user.errors.full_messages }, status: :unprocessable_entity
      end

    else
    
      Rails.logger.info "Member updating self with id: #{current_user.id}"
      authorize current_user

    
      if current_user.member? && params[:user][:password].present?
        unless current_user.valid_password?(params[:user][:current_password])
          Rails.logger.warn "Current password incorrect for user #{current_user.id}"
          return render json: { message: 'Current password is incorrect' }, status: :unauthorized
        end
      end

      before_attrs = current_user.slice(:name, :email)

      if current_user.update(account_update_params)
        Rails.logger.info "User updated successfully: #{current_user.inspect}"

        after_attrs = current_user.slice(:name, :email)

        if current_user.member?
          MemberLog.create!(
            action: "update",
            member: {
              id: current_user.id,
              name: current_user.name,
              email: current_user.email
            },
            before_update: before_attrs,
            after_update: after_attrs
          )
        end

        render json: { message: 'User updated successfully', user: user_response(current_user) }, status: :ok
      else
        Rails.logger.info "Update failed for user: #{current_user.errors.full_messages}"
        render json: { message: 'Update failed', errors: current_user.errors.full_messages }, status: :unprocessable_entity
      end
    end
  end


  def destroy
    Rails.logger.info "Destroy action started by user #{current_user&.id} with params: #{params.inspect}"

    if current_user.admin? && params[:id].present?
      target_user = User.find_by(id: params[:id], role: :librarian)

      unless target_user
        Rails.logger.info "Librarian not found for deletion with id: #{params[:id]}"
        return render json: { message: 'Librarian not found' }, status: :not_found
      end

      authorize target_user

      librarian_data = target_user.slice(:id, :name, :email)

      if target_user.destroy
        Rails.logger.info "Librarian deleted successfully: #{target_user.inspect}"

        LibrarianLog.create!(
          action: "delete",
          admin: {
            id: current_user.id,
            name: current_user.name,
            email: current_user.email
          },
          librarian: librarian_data
        )

        render json: { message: 'Librarian deleted successfully' }, status: :ok
      else
        Rails.logger.info "Failed to delete librarian: #{target_user.errors.full_messages}"
        render json: { message: 'Failed to delete librarian', errors: target_user.errors.full_messages }, status: :unprocessable_entity
      end
    else
      target_user = current_user
      authorize target_user

      member_data = target_user.slice(:id, :name, :email)

      if target_user.destroy
        Rails.logger.info "Account deleted successfully: #{target_user.inspect}"

        if target_user.member?
          MemberLog.create!(
            action: "delete",
            member: member_data
          )
        end

        render json: { message: 'Account deleted successfully' }, status: :ok
      else
        Rails.logger.info "Failed to delete account: #{target_user.errors.full_messages}"
        render json: { message: 'Account deletion failed', errors: target_user.errors.full_messages }, status: :unprocessable_entity
      end
    end
  end

  private

  def sign_up_params
    params.require(:user).permit(:email, :password, :password_confirmation, :name, :role)
  end

  def account_update_params
    params.require(:user).permit(:email, :password, :password_confirmation, :name)
  end

  def user_response(user)
    return {} unless user.present?

    {
      id: user.id,
      email: user.email,
      name: user.name,
      role: user.role,
      created_at: user.created_at
    }
  end
end
