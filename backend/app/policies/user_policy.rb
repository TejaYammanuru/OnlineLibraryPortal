class UserPolicy
  attr_reader :current_user, :record

  def initialize(current_user, record)
    @current_user = current_user
    @record = record
  end

  def create?
    return record.role == 'member' if current_user.nil? 

    current_user.admin? && record.role == 'librarian'
  end

  def update?
    return false unless current_user.present?

    
    return true if current_user.admin? && record.librarian?

    
    current_user == record && current_user.member?
  end

  def destroy?
    return false unless current_user.present?

   
    return true if current_user.admin? && record.librarian?

    
    current_user == record && current_user.member?
  end
end
